package main

import (
	"chym/stream/backend/api"
	"chym/stream/backend/config"
	"chym/stream/backend/db"
	"chym/stream/backend/utils"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
)

type Option struct {
	WorkPath string `short:"w" long:"work" description:"媒体库路径" default:"E:\\media"`
}

var opt Option

func main() {
	_, err := flags.Parse(&opt)
	if err != nil {
		log.Println("参数解析错误")
		return
	}

	if opt.WorkPath == "" {
		log.Println("未设置媒体库路径")
		return
	}

	config.InitConfig(config.AppConfig{
		WorkPath: opt.WorkPath,
	})

	db.InitDB()
	utils.InitDowner(3)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Upgrade", "Origin", "Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.POST("/import/media", api.ImportMediaHandler)
	r.POST("/list/media", api.ListHandler)
	r.POST("/get/media", api.GetMediaHandler)
	r.POST("/down/media", api.DownloadMediaHandle)
	r.POST("/list/task", api.DownTaskListHandle)

	r.Run("0.0.0.0:8080")
}
