package utils

import (
	"chym/stream/backend/db"
	"testing"
)

// func TestReadAllFromUrl(t *testing.T) {
// 	b, err := ReadAllFromUrl("https://v.cdnlz1.com/20230911/22053_33fb8c1e/2000k/hls/mixed.m3u8")
// 	t.Log(err)
// 	strs := strings.Split(string(b), "\n")
// 	OutputNewM3u8(strs, "https://v.cdnlz1.com/20230911/22053_33fb8c1e/2000k/hls", "E:\\media\\tvs\\1.m3u8")
// 	if err != nil {
// 		t.Error(`TestReadAllFromUrl = false`)
// 	}
// }

func TestParseTvShowXml(t *testing.T) {
	media := db.Media{}
	ParseTvShowXml("E:\\media\\tvs\\异人之下(2023)\\tvshow.nfo", &media)
	t.Log(media)
}
