package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/tidwall/sjson"
)

type SetupConfig struct {
	MediaPath string `json:"media_path" bson:"media_path"`
	DBPath    string `json:"db_path" bson:"db_path"`
	TmmPath   string `json:"tmm_path" bson:"tmm_path"`
}

type AppConfig struct {
	MediaPath string
	DBPath    string
	MoviePath string
	TvPath    string
	TmmPath   string
}

var AppConf = AppConfig{
	MediaPath: "",
}

func InitConfig(configPath string) AppConfig {
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

	if setupConfig.TmmPath == "" {
		panic("缺少TMM路径")
	}

	AppConf.MediaPath = setupConfig.MediaPath
	AppConf.DBPath = setupConfig.DBPath
	AppConf.TmmPath = setupConfig.TmmPath

	if _, err := os.Stat(AppConf.TmmPath); os.IsNotExist(err) {
		panic("TMM路径不存在")
	}

	AppConf.MoviePath = filepath.Join(AppConf.MediaPath, "movies")
	AppConf.TvPath = filepath.Join(AppConf.MediaPath, "tvs")

	_ = os.Mkdir(AppConf.MoviePath, os.ModeDir)
	_ = os.Mkdir(AppConf.TvPath, os.ModeDir)

	updateTmmMediaPath()

	log.Println(AppConf.MoviePath)
	log.Println(AppConf.TvPath)

	return AppConf
}

func updateTmmMediaPath() {
	tmmTvConfigPath := filepath.Join(AppConf.TmmPath, "data", "tvShows.json")
	tmmMovieConfigPath := filepath.Join(AppConf.TmmPath, "data", "movies.json")

	tvMaps := make(map[string]interface{})
	tvMaps["tvShowDataSource"] = [1]string{AppConf.TvPath}
	err := UpdateJsonKey(tmmTvConfigPath, tvMaps, true)
	if err != nil {
		panic("设置TMM出错")
	}

	movieMaps := make(map[string]interface{})
	movieMaps["movieDataSource"] = [1]string{AppConf.MoviePath}
	err = UpdateJsonKey(tmmMovieConfigPath, movieMaps, true)
	if err != nil {
		panic("设置TMM出错")
	}
}

func UpdateJsonKey(jsonPath string, maps map[string]interface{}, newBack bool) error {
	b, err := os.ReadFile(jsonPath) // just pass the file name
	if err != nil {
		return err
	}
	if newBack {
		err = ioutil.WriteFile(jsonPath+".back", b, 0644)
		if err != nil {
			return err
		}
	}
	content := string(b[:])
	for k, v := range maps {
		content, _ = sjson.Set(content, k, v)
	}
	newContent := []byte(content)
	err = os.WriteFile(jsonPath, newContent, 0644)
	if err != nil {
		return err
	}
	return nil
}
