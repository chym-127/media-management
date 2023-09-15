package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Media{}, &MediaWithTag{}, &MediaWithActor{}, &Actor{}, &MediaTag{})
}
