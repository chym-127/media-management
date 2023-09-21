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
	record, _ := db.GetMediaDownloadRecordByMediaID(media.ID)
	if record.Type == 1 || record.Type == 2 {
		c.JSON(http.StatusOK, GenResponse(nil, TASK_RUNNING, "下载任务已存在，请勿重复下载"))
		return
	}
	go utils.DownloadMediaAllEpisode(media, record)
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
			SuccessCount:  v.SuccessCount,
			FailedCount:   v.FailedCount,
			Type:          v.Type,
		}
		resp = append(resp, m)
	}

	c.JSON(http.StatusOK, GenResponse(resp, SUCCESS, "SUCCESS"))
}
