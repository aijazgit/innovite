package store

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"example.com/monitor"
)

type IDB interface {
	Create(string, int64, *zap.SugaredLogger) error
	Update(string, int64, *zap.SugaredLogger) error
	Save(chan monitor.FileInfo, *zap.SugaredLogger)
}

type DB struct {
	db *gorm.DB
}

type FileInfo struct {
	PathFileName string `gorm:"primaryKey;"`
	Size         int64
}
