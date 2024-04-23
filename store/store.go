package store

import (
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(dbName string, slog *zap.SugaredLogger) (IDB, error) {
	handle, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		slog.Errorw("failed to open database",
			"dbName", dbName,
			"error", err)
		return nil, err
	}

	handle.AutoMigrate(&FileInfo{})

	return &DB{handle}, nil
}

func (d *DB) Create(pathFileName string, size int64, slog *zap.SugaredLogger) error {
	db := d.db

	insertFileInfo := &FileInfo{PathFileName: pathFileName, Size: size}

	result := db.Create(insertFileInfo)
	if result.Error != nil {
		slog.Errorw("failed to create record",
			"pathFileName", pathFileName,
			"size", size,
			"error", result.Error)
		return result.Error
	}

	return nil
}

func (d *DB) Update(pathFileName string, size int64, slog *zap.SugaredLogger) error {
	db := d.db

	updateFileInfo := &FileInfo{PathFileName: pathFileName, Size: size}

	result := db.Save(updateFileInfo)
	if result.Error != nil {
		slog.Errorw("failed to update record",
			"pathFileName", pathFileName,
			"size", size,
			"error", result.Error)
		return result.Error
	}

	return nil
}
