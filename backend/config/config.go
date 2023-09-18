package config

import (
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	WorkPath  string `json:"workPath"`
	MoviePath string `json:"moviePath"`
	TvPath    string `json:"tvPath"`
}

var AppConf = AppConfig{
	WorkPath: "",
}

func InitConfig(config AppConfig) {
	AppConf.WorkPath = config.WorkPath
	if _, err := os.Stat(AppConf.WorkPath); os.IsNotExist(err) {
		log.Println("目录不存在")
		panic("目录不存在")
	}

	if AppConf.WorkPath != "" {
		AppConf.MoviePath = filepath.Join(AppConf.WorkPath, "movies")
		AppConf.TvPath = filepath.Join(AppConf.WorkPath, "tvs")
		_ = os.Mkdir(AppConf.MoviePath, os.ModeDir)
		_ = os.Mkdir(AppConf.TvPath, os.ModeDir)
		log.Println(AppConf.MoviePath)
		log.Println(AppConf.TvPath)

	}
}
