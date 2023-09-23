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

	_, err := db.ListMedia(protocols.ListMediaReq{})
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
	utils.GetMediaMetaFromTMDB(1)
	utils.GetMediaMetaFromTMDB(2)
	args := protocols.ListMediaReq{
		Page: protocols.Page{
			PageLimit: -1,
		},
	}
	medias, err := db.ListMedia(args)
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
		count := 0
		if m.Type == 2 {
			mediaPath := filepath.Join(config.AppConf.TvPath, m.Title+"("+strconv.Itoa(int(m.ReleaseDate))+")")
			xmlFilePath := filepath.Join(mediaPath, "tvshow.nfo")
			utils.ParseTvShowXml(xmlFilePath, &m)
			count, _ = utils.UpdateTvShowEpisodeFileName(episode, mediaPath)
			utils.ParseTvShowEpisodeXml(episode, &m, mediaPath)
		}
		if m.Type == 1 {
			mediaPath := filepath.Join(config.AppConf.MoviePath, m.Title+"("+strconv.Itoa(int(m.ReleaseDate))+")")
			xmlFilePath := filepath.Join(mediaPath, "movie.nfo")
			count, _ = utils.UpdateMovieEpisodeFileName(mediaPath, m.Title)
			utils.ParseMovieXml(xmlFilePath, &m, episode)
		}
		m.LocalEpisodeCount = int16(count)
		db.UpdateMedia(&m)
	}

	c.JSON(http.StatusOK, GenResponse(nil, SUCCESS, "SUCCESS"))
}
