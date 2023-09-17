package utils

import (
	"strings"
	"testing"
)

func TestReadAllFromUrl(t *testing.T) {
	b, err := ReadAllFromUrl("https://v.cdnlz1.com/20230911/22053_33fb8c1e/2000k/hls/mixed.m3u8")
	t.Log(err)
	strs := strings.Split(string(b), "\n")
	OutputNewM3u8(strs, "https://v.cdnlz1.com/20230911/22053_33fb8c1e/2000k/hls", "C:\\Medias\\tv\\1.m3u8")
	if err != nil {
		t.Error(`TestReadAllFromUrl = false`)
	}

}
