package database

import (
	"github.com/evellyncosta/context-go/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Cotacao{}); err != nil {
		return nil, err
	}

	return db, nil
}