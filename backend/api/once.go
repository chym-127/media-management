package api

import (
	"chym/stream/backend/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var UpdateDiskFromDBRunning = false

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
