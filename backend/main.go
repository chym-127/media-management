package main

import (
	"chym/stream/backend/api"
	"chym/stream/backend/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()
	r.GET("/ping", api.ImportMediaHandler)
	r.Run()
}
