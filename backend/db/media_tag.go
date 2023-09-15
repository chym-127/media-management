package db

import "gorm.io/gorm"

type MediaTag struct {
	gorm.Model
	Name        string
	Description string
}
