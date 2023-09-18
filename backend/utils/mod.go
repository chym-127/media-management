package utils

import (
	"chym/stream/backend/config"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func ReadAllFromUrl(url string) (b []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func RemoveFileIfExist(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		err = os.Remove(filePath) //remove the file using built-in functions
		if err == nil {
			return true
		} else {
			log.Println(err)
		}
	}
	return false
}

func GetMediaPath(mediaType int8) string {
	// 电影
	if mediaType == 1 {
		return filepath.Join(config.AppConf.WorkPath, "movies")
	}

	// 电视剧
	if mediaType == 2 {
		return filepath.Join(config.AppConf.WorkPath, "tvs")
	}

	return ""
}
