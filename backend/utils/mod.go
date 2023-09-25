package utils

import (
	"chym/stream/backend/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/sjson"
)

func ReadAllFromUrl(url string) (b []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	log.Println(req)
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
		return config.AppConf.MoviePath
	}

	// 电视剧
	if mediaType == 2 {
		return config.AppConf.TvPath
	}

	return ""
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
