package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type SetupConfig struct {
	MediaPath string `json:"media_path" bson:"media_path"`
	DBPath    string `json:"db_path" bson:"db_path"`
}

type AppConfig struct {
	MediaPath string
	DBPath    string
	MoviePath string
	TvPath    string
}

var AppConf = AppConfig{
	MediaPath: "",
}

func InitConfig(configPath string) {
	b, err := os.ReadFile(configPath) // just pass the file name
	if err != nil {
		panic("配置文件不存在")
	}
	setupConfig := SetupConfig{}
	err = json.Unmarshal(b, &setupConfig)
	if err != nil {
		panic("配置文件解析失败")
	}

	if setupConfig.MediaPath == "" {
		panic("缺少媒体路径")
	}

	if setupConfig.DBPath == "" {
		panic("缺少数据库路径")
	}

	AppConf.MediaPath = setupConfig.MediaPath
	AppConf.DBPath = setupConfig.DBPath

	AppConf.MoviePath = filepath.Join(AppConf.MediaPath, "movies")
	AppConf.TvPath = filepath.Join(AppConf.MediaPath, "tvs")
	_ = os.Mkdir(AppConf.MoviePath, os.ModeDir)
	_ = os.Mkdir(AppConf.TvPath, os.ModeDir)
	log.Println(AppConf.MoviePath)
	log.Println(AppConf.TvPath)
}
