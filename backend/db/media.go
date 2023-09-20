package db

import (
	"chym/stream/backend/protocols"
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	Title       string
	ReleaseDate int16
	Description string
	Score       float64
	Episodes    string
	PlayConfig  string
	PosterUrl   string
	FanartUrl   string
	Area        string
	Type        int8
	Expand      string
}

type MediaDownloadRecord struct {
	gorm.Model
	Title         string
	MediaID       uint
	EpisodeCount  uint
	DownloadCount uint
	Type          uint //1队列中 2下载中 3下载成功 4下载失败
}

type MediaWithTag struct {
	gorm.Model
	MediaID uint
	TagID   uint
}

type MediaWithActor struct {
	gorm.Model
	MediaID uint
	ActorID uint
}

func CreateMedia(mediaItem protocols.MediaItem) (Media, error) {
	jsonByte, _ := json.Marshal(mediaItem.Episodes)
	media := Media{
		Title:       mediaItem.Title,
		ReleaseDate: mediaItem.ReleaseDate,
		Episodes:    string(jsonByte[:]),
		Type:        mediaItem.Type,
	}

	DB.Create(&media)

	return media, nil
}

func UpdateMedia(media *Media) error {
	DB.Save(media)
	return nil
}

func ListMedia() ([]Media, error) {
	var medias []Media
	result := DB.Model(&Media{}).Find(&medias)
	if result.Error != nil {
		return medias, result.Error
	}
	return medias, nil
}

func GetMediaByID(id uint) (Media, error) {
	var media Media
	result := DB.First(&media, id)
	if result.Error != nil {
		log.Println(result.Error)
		return media, result.Error
	}
	return media, nil
}

func GetMediaByTitleWithDate(title string, date int16) (Media, error) {
	var media Media
	result := DB.Where("title = ?", title).Where("release_date = ?", date).First(&media)
	if result.Error != nil {
		return media, result.Error
	}
	return media, nil
}

func CreateMediaDownRecord(record MediaDownloadRecord) (MediaDownloadRecord, error) {
	DB.Create(&record)
	return record, nil
}

func UpdateMediaDownRecord(record *MediaDownloadRecord) error {
	DB.Save(record)
	return nil
}

func GetMediaDownloadRecordByMediaID(id uint) (MediaDownloadRecord, error) {
	var mediaDownloadRecord MediaDownloadRecord
	result := DB.Where("media_id = ?", id).First(&mediaDownloadRecord)
	if result.Error != nil {
		return mediaDownloadRecord, result.Error
	}
	return mediaDownloadRecord, nil
}

func ListMediaDownloadRecord() ([]MediaDownloadRecord, error) {
	var mediaDownloadRecords []MediaDownloadRecord
	result := DB.Model(&MediaDownloadRecord{}).Find(&mediaDownloadRecords)
	if result.Error != nil {
		log.Println(result.Error)
	}
	return mediaDownloadRecords, nil
}
