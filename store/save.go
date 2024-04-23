package store

import (
	"os"

	"example.com/monitor"
	"go.uber.org/zap"
)

func (d *DB) Save(fileChan chan monitor.FileInfo, slog *zap.SugaredLogger) {
	slog.Infow("loop starting ...")
	defer slog.Infow("loop leaving ...")

	for fileInfo := range fileChan {
		info, err := os.Stat(fileInfo.PathFileName)
		if err != nil {
			slog.Errorw("failed to get file os.Stat()",
				"pathFileName", fileInfo.PathFileName,
				"error", err)
			return
		}

		if fileInfo.IsCreated {
			slog.Infow("file created", "pathFileName", fileInfo.PathFileName, "size", info.Size())
			d.Create(fileInfo.PathFileName, info.Size(), slog)
		} else if fileInfo.IsModified {
			slog.Infow("file updated", "pathFileName", fileInfo.PathFileName, "size", info.Size())
			d.Update(fileInfo.PathFileName, info.Size(), slog)
		}
	}
}
