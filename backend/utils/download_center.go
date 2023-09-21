package utils

import (
	"chym/stream/backend/config"
	"chym/stream/backend/db"
	"chym/stream/backend/protocols"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/panjf2000/ants/v2"
)

var DOWNER *ants.Pool

// 初始化下载器
func InitDowner(max int) {
	var err error
	DOWNER, err = ants.NewPool(max, ants.WithOptions(ants.Options{
		MaxBlockingTasks: 1000,
	}))
	if err != nil {
		panic("InitDowner failed")
	}
}

func taskFuncWrapper(m3u8FilePath string, outputFilePath string, mediaId uint) taskFunc {
	return func() {
		args := []string{"-i", m3u8FilePath, "-o", outputFilePath, "-c", "-g"}
		log.Println("m3u8_downloader command :m3u8_downloader " + strings.Join(args, " "))
		cmd := exec.Command("m3u8_downloader", args...)
		cmd.Dir = filepath.Dir(m3u8FilePath)
		_, err := cmd.CombinedOutput()
		success := true
		if err != nil {
			success = false
		}

		record, err := db.GetMediaDownloadRecordByMediaID(mediaId)
		if err == nil {
			record.DownloadCount += 1
			if success {
				record.SuccessCount += 1
			} else {
				record.FailedCount += 1
			}
			if record.DownloadCount == record.EpisodeCount {
				record.Type = 3
			} else {
				record.Type = 2
			}
			db.UpdateMediaDownRecord(&record)
		}
	}
}

type DownTask struct {
	MeidaID        uint
	OutputFilePath string
	InputFileName  string
	NfoFileName    string
}

func DownloadMediaAllEpisode(media db.Media) error {
	var mediaPath string
	if media.Type == 2 {
		mediaPath = filepath.Join(config.AppConf.TvPath, media.Title+"("+strconv.Itoa(int(media.ReleaseDate))+")")
	}
	if media.Type == 1 {
		mediaPath = filepath.Join(config.AppConf.MoviePath, media.Title+"("+strconv.Itoa(int(media.ReleaseDate))+")")
	}
	var episodes []protocols.EpisodeItem
	var tasks []DownTask
	err := json.Unmarshal([]byte(media.Episodes), &episodes)
	if err != nil {
		return err
	}
	record := db.MediaDownloadRecord{
		Title:         media.Title,
		MediaID:       media.ID,
		SuccessCount:  0,
		FailedCount:   0,
		DownloadCount: 0,
		EpisodeCount:  uint(len(episodes)),
		Type:          1,
	}

	str := "Season-"
	for _, episode := range episodes {
		var inputFileName string
		var outputFileName string
		var nfoFileName string
		seasonName := str + strconv.Itoa(int(episode.Season))
		if media.Type == 1 {
			inputFileName = media.Title + ".m3u8"
			outputFileName = media.Title + ".mp4"
			nfoFileName = "movie"
		}
		if media.Type == 2 {
			inputFileName = filepath.Join(seasonName, "E"+strconv.Itoa(int(episode.Index))+".m3u8")
			outputFileName = filepath.Join(seasonName, "E"+strconv.Itoa(int(episode.Index))+".mp4")
			nfoFileName = "E" + strconv.Itoa(int(episode.Index))
		}
		m3u8FilePath := filepath.Join(mediaPath, inputFileName)
		outputFilePath := filepath.Join(mediaPath, outputFileName)
		if _, err := os.Stat(outputFilePath); errors.Is(err, os.ErrNotExist) {
			tasks = append(tasks, DownTask{
				InputFileName:  m3u8FilePath,
				OutputFilePath: outputFilePath,
				MeidaID:        media.ID,
				NfoFileName:    nfoFileName,
			})
		} else {
			record.DownloadCount += 1
			record.SuccessCount += 1
		}
	}
	if record.DownloadCount == uint(len(episodes)) {
		record.Type = 3
	}
	db.CreateMediaDownRecord(record)
	for _, task := range tasks {
		DOWNER.Submit(taskFuncWrapper(task.InputFileName, task.OutputFilePath, task.MeidaID))
	}
	return nil
}
