package db

import (
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
