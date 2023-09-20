package api

import (
	"chym/stream/backend/db"
	"chym/stream/backend/protocols"
	"chym/stream/backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DownloadMediaHandle(c *gin.Context) {
	getMediaReq := protocols.GetMediaReq{}
	err := c.ShouldBindJSON(&getMediaReq)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, GenResponse(nil, PARAMETER_ERROR, "FAILED"))
		return
	}
	media, err := db.GetMediaByID(getMediaReq.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, GenResponse(nil, FAILED, "FAILED"))
		return
	}
	_, err = db.GetMediaDownloadRecordByMediaID(media.ID)
	if err == nil {
		c.JSON(http.StatusOK, GenResponse(nil, TASK_RUNNING, "下载任务已存在，请勿重复下载"))
		return
	}
	utils.DownloadMediaAllEpisode(media)
	c.JSON(http.StatusOK, GenResponse(nil, SUCCESS, "SUCCESS"))
}

func DownTaskListHandle(c *gin.Context) {
	items, err := db.ListMediaDownloadRecord()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, GenResponse(nil, FAILED, "FAILED"))
		return
	}
	var resp []protocols.MediaDownloadRecordItem
	for _, v := range items {
		m := protocols.MediaDownloadRecordItem{
			ID:            v.ID,
			Title:         v.Title,
			DownloadCount: v.DownloadCount,
			EpisodeCount:  v.EpisodeCount,
			Type:          v.Type,
		}
		resp = append(resp, m)
	}

	c.JSON(http.StatusOK, GenResponse(resp, SUCCESS, "SUCCESS"))
}
