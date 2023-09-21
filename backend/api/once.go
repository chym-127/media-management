package api

import (
	"chym/stream/backend/config"
	"chym/stream/backend/db"
	"chym/stream/backend/protocols"
	"chym/stream/backend/utils"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

var UpdateDiskFromDBRunning = false
var UpdateMediaMetaDataFromDiskRunning = false

func UpdateDiskFromDB(c *gin.Context) {
	if UpdateDiskFromDBRunning {
		c.JSON(http.StatusOK, GenResponse(nil, TASK_RUNNING, "任务进行中"))
		return
	}

	_, err := db.ListMedia()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, GenResponse(nil, FAILED, "FAILED"))
	}
}

func UpdateMediaMetaDataFromDisk(c *gin.Context) {
	if UpdateMediaMetaDataFromDiskRunning {
		c.JSON(http.StatusOK, GenResponse(nil, TASK_RUNNING, "任务进行中"))
		return
	}

	medias, err := db.ListMedia()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, GenResponse(nil, FAILED, "FAILED"))
		return
	}

	for _, m := range medias {
		var episode []protocols.EpisodeItem
		err = json.Unmarshal([]byte(m.Episodes), &episode)
		if err != nil {
			log.Println(err)
			continue
		}
		if m.Type == 2 {
			mediaPath := filepath.Join(config.AppConf.TvPath, m.Title+"("+strconv.Itoa(int(m.ReleaseDate))+")")
			xmlFilePath := filepath.Join(mediaPath, "tvshow.nfo")
			utils.ParseTvShowXml(xmlFilePath, &m)
			utils.UpdateTvShowEpisodeFileName(episode, mediaPath)
			utils.ParseTvShowEpisodeXml(episode, &m, mediaPath)
		}
		if m.Type == 1 {
			mediaPath := filepath.Join(config.AppConf.MoviePath, m.Title+"("+strconv.Itoa(int(m.ReleaseDate))+")")
			xmlFilePath := filepath.Join(mediaPath, "movie.nfo")
			utils.UpdateMovieEpisodeFileName(mediaPath, m.Title)
			utils.ParseMovieXml(xmlFilePath, &m, episode)
		}
		db.UpdateMedia(&m)
	}

	c.JSON(http.StatusOK, GenResponse(nil, SUCCESS, "SUCCESS"))
}
