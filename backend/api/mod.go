package api

import (
	"chym/stream/backend/db"
	"chym/stream/backend/protocols"
	"chym/stream/backend/utils"
	"log"
	"net/http"

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
	for _, media := range importMediaReqProtocol.Medias {
		db.CreateMedia(media)
		utils.SaveEpisode2Disk("C://Medias//tv", media.Episodes)
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
