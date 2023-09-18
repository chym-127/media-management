package utils

import (
	"bufio"
	"chym/stream/backend/config"
	"chym/stream/backend/db"
	"chym/stream/backend/protocols"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beevik/etree"
	"github.com/panjf2000/ants/v2"
)

func GetBaseUrl(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		log.Println(err)
	}
	baseUrl := u.Scheme + "://" + u.Host
	fileNameReg, _ := regexp.Compile(".?.m3u8")
	name := path.Base(uri)
	if fileNameReg.MatchString(name) {
		baseUrl += path.Dir(u.Path)
	} else {
		baseUrl += u.Path
	}
	return baseUrl
}

func OutputNewM3u8(strArr []string, baseUrl string, outputPath string) error {
	url_regexp, _ := regexp.Compile("http.*.ts")
	ts_regexp, _ := regexp.Compile(".*.ts$")
	duration_regexp, _ := regexp.Compile("#EXTINF:.*")
	number_re := regexp.MustCompile("[0-9]+.[0-9]+")
	number_re1 := regexp.MustCompile("[0-9]+")

	outFile, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
		return err
	}

	if err := os.Truncate(outputPath, 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
		return err
	}

	datawriter := bufio.NewWriter(outFile)
	if err != nil {
		return err
	}

	for _, line := range strArr {
		newLine := line
		if ts_regexp.MatchString(line) {
			if !url_regexp.MatchString(line) {
				newLine = baseUrl + "/" + line
			}
		}
		if duration_regexp.MatchString(line) {
			str := number_re.FindString(line)
			if len(str) == 0 {
				str = number_re1.FindString(line)
			}
			if len(str) >= 4 {
				str = str[:4]
			}
			newLine = "#EXTINF: " + str + ","
		}
		_, _ = datawriter.WriteString(newLine + "\n")
	}
	datawriter.Flush()
	outFile.Close()

	return nil
}

var RETRY_DOWNLOAD_COUNT = 3

type taskFunc func()

func TaskFuncWrapper(url string, outputPath string, wg *sync.WaitGroup) taskFunc {
	return func() {
		baseUrl := GetBaseUrl(url)
		var b []byte
		var err error
		for i := 1; i <= RETRY_DOWNLOAD_COUNT; i++ {
			b, err = ReadAllFromUrl(url)
			time.Sleep(time.Millisecond * time.Duration(300))
			if err == nil {
				break
			}
		}

		strs := strings.Split(string(b), "\n")
		err = OutputNewM3u8(strs, baseUrl, outputPath)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}
}

// 保存媒体资源的m3u8文件
func SaveEpisode2Disk(mediaItem protocols.MediaItem) {
	var wg sync.WaitGroup
	p, _ := ants.NewPool(5)
	defer p.Release()
	var fileName string
	var workPath string
	if mediaItem.Type == 1 {
		workPath = filepath.Join(config.AppConf.MoviePath, mediaItem.Title+"("+strconv.Itoa(int(mediaItem.ReleaseDate))+")")
	} else if mediaItem.Type == 2 {
		workPath = filepath.Join(config.AppConf.TvPath, mediaItem.Title+"("+strconv.Itoa(int(mediaItem.ReleaseDate))+")")
	}
	_ = os.Mkdir(workPath, os.ModeDir)

	//根据季数排序
	sort.SliceStable(mediaItem.Episodes, func(i, j int) bool {
		return mediaItem.Episodes[i].Season < mediaItem.Episodes[j].Season
	})

	for _, item := range mediaItem.Episodes {
		outputPath := workPath
		if mediaItem.Type == 1 {
			fileName = mediaItem.Title + ".m3u8"
		} else if mediaItem.Type == 2 {
			outputPath = filepath.Join(workPath, "Season "+strconv.Itoa(int(item.Season)))
			_ = os.Mkdir(outputPath, os.ModeDir)
			fileName = "E" + strconv.Itoa(int(item.Index)) + ".m3u8"
		}
		filePath := filepath.Join(outputPath, fileName)
		url := item.Url

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			wg.Add(1)
			p.Submit(TaskFuncWrapper(url, filePath, &wg))
		}
	}

	wg.Wait()

}

// 调用 tmm 进行刮削
func GetMediaMetaFromTMDB(mediaType int8) {
	var t string
	if mediaType == 1 {
		t = "movie"
	} else {
		t = "tvshow"
	}
	args := []string{t, "-u", "-n"}
	log.Println("tinyMediaManager command :tinyMediaManagerCMD " + strings.Join(args, " "))
	cmd := exec.Command("tinyMediaManagerCMD", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		panic("")
	}
}

func ParseTvShowXml(filePath string, mediaModal *db.Media) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filePath); err != nil {
		return err
	}
	mediaModal.Description = doc.FindElement("//tvshow/plot").Text()
	mediaModal.PosterUrl = doc.FindElement("//tvshow/thumb/[@aspect='poster']").Text()
	if doc.FindElement("//tvshow/fanart") != nil && doc.FindElement("//tvshow/fanart/thumb") != nil {
		mediaModal.FanartUrl = doc.FindElement("//tvshow/fanart/thumb").Text()
	}
	if doc.FindElement("//tvshow/ratings/rating") != nil {
		val, _ := strconv.ParseFloat(doc.FindElement("//tvshow/ratings/rating/[@name='themoviedb']/value").Text(), 32)
		mediaModal.Score = val
	}
	mediaModal.Area = doc.FindElement("//tvshow/country").Text()

	return nil
}

func ParseMovieXml(filePath string, mediaModal *db.Media) error {
	log.Println("ParseMovieXml")
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filePath); err != nil {
		log.Println(err)
		return err
	}

	mediaModal.Description = doc.FindElement("//movie/plot").Text()
	mediaModal.PosterUrl = doc.FindElement("//movie/thumb/[@aspect='poster']").Text()
	if doc.FindElement("//movie/fanart") != nil && doc.FindElement("//movie/fanart/thumb") != nil {
		mediaModal.FanartUrl = doc.FindElement("//movie/fanart/thumb").Text()
	}
	if doc.FindElement("//movie/ratings/rating") != nil {
		val, _ := strconv.ParseFloat(doc.FindElement("//movie/ratings/rating/[@name='themoviedb']/value").Text(), 32)
		mediaModal.Score = val
	}
	mediaModal.Area = doc.FindElement("//movie/country").Text()

	return nil
}