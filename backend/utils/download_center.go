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
	DOWNER, err = ants.NewPool(max)
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
		if err != nil {
			log.Println(err)
		}

		record, err := db.GetMediaDownloadRecordByMediaID(mediaId)
		if err == nil {
			record.DownloadCount += 1
			if record.DownloadCount == record.EpisodeCount {
				record.Type = 3
			} else {
				record.Type = 2
			}
			db.UpdateMediaDownRecord(&record)
		}
	}
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
	err := json.Unmarshal([]byte(media.Episodes), &episodes)
	if err != nil {
		return err
	}
	record := db.MediaDownloadRecord{
		Title:         media.Title,
		MediaID:       media.ID,
		DownloadCount: 0,
		EpisodeCount:  uint(len(episodes)),
		Type:          1,
	}

	db.CreateMediaDownRecord(record)
	str := "Season-"
	for _, episode := range episodes {
		var inputFileName string
		var outputFileName string
		seasonName := str + strconv.Itoa(int(episode.Season))
		if media.Type == 1 {
			inputFileName = media.Title + ".m3u8"
			outputFileName = media.Title + ".mp4"
		}
		if media.Type == 2 {
			inputFileName = filepath.Join(seasonName, "E"+strconv.Itoa(int(episode.Index))+".m3u8")
			outputFileName = filepath.Join(seasonName, "E"+strconv.Itoa(int(episode.Index))+".mp4")
		}
		m3u8FilePath := filepath.Join(mediaPath, inputFileName)
		outputFilePath := filepath.Join(mediaPath, outputFileName)
		if _, err := os.Stat(outputFilePath); errors.Is(err, os.ErrNotExist) {
			DOWNER.Submit(taskFuncWrapper(m3u8FilePath, outputFilePath, media.ID))
		}
	}

	return nil
}
