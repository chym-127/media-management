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

type RespCode int64

const (
	SUCCESS         RespCode = 200 //成功
	PARAMETER_ERROR RespCode = 201 //参数解析错误
	FAILED          RespCode = 202 //操作失败
)

type BaseResponse struct {
	Data    interface{} `json:"data"`
	Code    RespCode    `json:"code"`
	Message string      `json:"message"`
}

func GenResponse(data interface{}, code RespCode, msg string) BaseResponse {
	if msg == "" {
		msg = "SUCCESS"
	}

	return BaseResponse{
		Data:    data,
		Code:    code,
		Message: msg,
	}
}

func ImportMediaHandler(c *gin.Context) {
	importMediaReqProtocol := protocols.ImportMediaReqProtocol{}
	err := c.ShouldBindJSON(&importMediaReqProtocol)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, GenResponse(nil, PARAMETER_ERROR, "FAILED"))
		return
	}
	var mediaModels []db.Media
	for _, media := range importMediaReqProtocol.Medias {
		if media.Type == 0 {
			if len(media.Episodes) > 1 {
				media.Type = 2
			} else {
				media.Type = 1
			}
		}
		mediaModel, _ := db.CreateMedia(media)
		mediaModels = append(mediaModels, mediaModel)
		// utils.SaveEpisode2Disk(media)
	}
	// utils.GetMediaMetaFromTMDB(1)
	// utils.GetMediaMetaFromTMDB(2)

	for _, m := range mediaModels {
		if m.Type == 2 {
			xmlFilePath := filepath.Join(config.AppConf.TvPath, m.Title+"("+strconv.Itoa(int(m.ReleaseDate))+")", "tvshow.nfo")
			utils.ParseTvShowXml(xmlFilePath, &m)
		}
		if m.Type == 1 {
			xmlFilePath := filepath.Join(config.AppConf.TvPath, m.Title+"("+strconv.Itoa(int(m.ReleaseDate))+")", "movie.nfo")
			utils.ParseMovieXml(xmlFilePath, &m)
		}
		db.UpdateMedia(&m)
	}

	c.JSON(http.StatusOK, GenResponse(nil, SUCCESS, "SUCCESS"))
}

func ListHandler(c *gin.Context) {
	listMediaReq := protocols.ListMediaReq{}
	err := c.ShouldBindJSON(&listMediaReq)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, GenResponse(nil, PARAMETER_ERROR, "FAILED"))
		return
	}
	medias, err := db.ListMedia()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, GenResponse(nil, FAILED, "FAILED"))
	}
	var resp []protocols.MediaItem
	for _, v := range medias {
		var episode []protocols.EpisodeItem
		err = json.Unmarshal([]byte(v.Episodes), &episode)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, GenResponse(nil, FAILED, "FAILED"))
		}
		m := protocols.MediaItem{
			Title:       v.Title,
			ReleaseDate: v.ReleaseDate,
			Description: v.Description,
			Score:       v.Score,
			Episodes:    episode,
			PlayConfig:  v.PlayConfig,
			PosterUrl:   v.PosterUrl,
			FanartUrl:   v.FanartUrl,
			Area:        v.Area,
			Type:        v.Type,
		}
		resp = append(resp, m)
	}

	c.JSON(http.StatusOK, GenResponse(resp, SUCCESS, "SUCCESS"))
}
