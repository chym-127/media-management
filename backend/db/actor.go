package db

import "gorm.io/gorm"

type Actor struct {
	gorm.Model
	Name        string
	Description string
	PosterUrl   string
	Expand      string
}
