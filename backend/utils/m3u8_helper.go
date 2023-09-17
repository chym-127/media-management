package utils

import (
	"bufio"
	"chym/stream/backend/protocols"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

type M3u8Info struct {
	name           string //名称
	is_aes         bool   // 是否加密
	key_str        string
	base_url       string
	total_segment  int32     //子元素的个数
	total_duration float64   //总时长
	segments       []Segment //子元素
}

type Segment struct {
	index        int32   //位置索引
	duration     float64 //时长
	duration_str string
	url          string //地址
	is_download  bool
}

func ParseM3u8(filePath string, name string, baseUrl string) M3u8Info {
	info := M3u8Info{
		name:           name,
		is_aes:         false,
		base_url:       baseUrl,
		total_segment:  0,
		total_duration: 0.0,
		segments:       []Segment{},
	}
	url_regexp, _ := regexp.Compile("http.*.ts")
	duration_regexp, _ := regexp.Compile("#EXTINF:.*")
	number_re := regexp.MustCompile("[0-9]+.[0-9]+")
	number_re1 := regexp.MustCompile("[0-9]+")

	readFile, err := os.Open(filePath)

	if err != nil {
		log.Println(err)
	} else {
		defer readFile.Close()
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var count_index int32 = 1
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if duration_regexp.MatchString(line) {
			str := number_re.FindString(line)
			if len(str) == 0 {
				str = number_re1.FindString(line)
			}
			var d float64 = 0.0
			if len(str) > 0 {
				v, _ := strconv.ParseFloat(str, 64)
				d, _ = strconv.ParseFloat(fmt.Sprintf("%.10f", v), 64)
			}
			s := Segment{
				index:        count_index,
				duration:     d,
				duration_str: str,
				is_download:  true,
			}
			if fileScanner.Scan() {
				line = fileScanner.Text()
				if url_regexp.MatchString(line) {
					s.url = line
				} else {
					s.url = info.base_url + "/" + line
				}
			}
			info.total_segment += 1
			info.total_duration += s.duration
			count_index += 1
			info.segments = append(info.segments, s)
		}
	}
	return info
}

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
		fmt.Println(err)
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
		log.Println(url)
		log.Println(outputPath)
		baseUrl := GetBaseUrl(url)
		var b []byte
		var err error
		for i := 1; i <= RETRY_DOWNLOAD_COUNT; i++ {
			b, err = ReadAllFromUrl(url)
			time.Sleep(time.Millisecond * time.Duration(500))
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

func SaveEpisode2Disk(workPath string, episodes []protocols.EpisodeItem) {
	var wg sync.WaitGroup
	p, _ := ants.NewPool(5)
	defer p.Release()

	for _, item := range episodes {
		wg.Add(1)
		filePath := filepath.Join(workPath, "E"+strconv.Itoa(int(item.Index))+".m3u8")
		url := item.Url
		p.Submit(TaskFuncWrapper(url, filePath, &wg))
	}

	wg.Wait()
}
