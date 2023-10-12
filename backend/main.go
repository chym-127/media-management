package main

import (
	"chym/stream/backend/api"
	"chym/stream/backend/config"
	"chym/stream/backend/db"
	"chym/stream/backend/iot"
	"chym/stream/backend/utils"
	"io"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
)

type Option struct {
	ConfigPath string `short:"c" long:"config" description:"配置文件"`
}

var opt Option

func main() {
	_, err := flags.Parse(&opt)
	if err != nil {
		log.Println("参数解析错误")
		return
	}

	if opt.ConfigPath == "" {
		log.Println("未设置配置文件")
		return
	}

	config.InitConfig(opt.ConfigPath)

	db.InitDB(config.AppConf.DBPath)
	utils.InitDowner(5)
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
	r.POST("/update/media", api.UpdateMediaHandler)
	r.POST("/list/media", api.ListHandler)
	r.POST("/get/media", api.GetMediaHandler)
	r.POST("/down/media", api.DownloadMediaHandle)
	r.POST("/list/task", api.DownTaskListHandle)
	r.POST("/down/media/nolocal", api.DownloadAllMediaHandle)

	r.POST("/update/medias/from-disk", api.UpdateMediaMetaDataFromDisk)
	r.POST("/update/medias/from-db", api.UpdateMediaLocalFromDB)

	api.NewServer()
	r.GET("/stream", api.MsgStream.ServeHTTP(), func(c *gin.Context) {
		v, ok := c.Get("clientChan")
		if !ok {
			return
		}
		clientChan, ok := v.(api.ClientChan)
		if !ok {
			return
		}
		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-clientChan; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 10)
	// 		now := time.Now().Format("2006-01-02 15:04:05")
	// 		currentTime := fmt.Sprintf("The Current Time Is %v", now)

	// 		// Send current time to clients message channel
	// 		api.MsgStream.Message <- currentTime
	// 	}
	// }()

	iot.InitSerial()

	r.Run("0.0.0.0:8080")
}
