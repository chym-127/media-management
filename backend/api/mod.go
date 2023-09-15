package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RespCode int64

const (
	SUCCESS         RespCode = 200 //成功
	PARAMETER_ERROR RespCode = 201 //参数解析错误
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

type ImportMediaReqProtocol struct {
	Medias []MediaItem `json:"medias"`
}

type MediaItem struct {
	Title       string        `json:"title" bson:"title" binding:"required"`
	ReleaseDate int16         `json:"releaseDate" bson:"release_date" binding:"required"`
	Description string        `json:"description" bson:"description"`
	Score       int16         `json:"score" bson:"score"`
	Episodes    []EpisodeItem `json:"episodes" bson:"episodes"`
	PlayConfig  string        `json:"playConfig" bson:"play_config"`
	PosterUrl   string        `json:"posterUrl" bson:"poster_url"`
	FanartUrl   string        `json:"fanartUrl" bson:"fanart_url"`
	Area        string        `json:"area" bson:"area"`
	Type        int8          `json:"type" bson:"type"`
}

type EpisodeItem struct {
	Url   string `json:"url" bson:"url"`
	Name  string `json:"name" bson:"name"`
	Index int8   `json:"index" bson:"index"`
}

func ImportMediaHandler(c *gin.Context) {
	importMediaReqProtocol := ImportMediaReqProtocol{}
	err := c.ShouldBindJSON(&importMediaReqProtocol)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, GenResponse(nil, PARAMETER_ERROR, "FAILED"))
		return
	}
	log.Println(importMediaReqProtocol)
	c.JSON(http.StatusOK, GenResponse(BaseResponse{}, SUCCESS, "SUCCESS"))
}
