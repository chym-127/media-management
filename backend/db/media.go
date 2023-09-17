package db

import (
	"chym/stream/backend/protocols"
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	Title       string
	ReleaseDate int16
	Description string
	Score       int16
	Episodes    string
	PlayConfig  string
	PosterUrl   string
	FanartUrl   string
	Area        string
	Type        int8
	Expand      string
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

func CreateMedia(mediaItem protocols.MediaItem) error {
	jsonByte, _ := json.Marshal(mediaItem.Episodes)
	media := Media{
		Title:       mediaItem.Title,
		ReleaseDate: mediaItem.ReleaseDate,
		Episodes:    string(jsonByte[:]),
	}

	result := DB.Create(&media)
	fmt.Println(result)

	return nil
}

func ListMedia() ([]Media, error) {
	var medias []Media
	result := DB.Model(&Media{}).Where("1=1").Find(&medias)
	if result.Error != nil {
		log.Println(result.Error)
	}
	return medias, nil
}
