package monitor

import (
	"errors"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

func GoMonitor(dirPath string, ch chan FileInfo, slog *zap.SugaredLogger) {
	slog.Infow("GoMonitor starting ...")

	var err error
	var watcher *fsnotify.Watcher

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}

			slog.Errorw("failed, monitor crashed",
				"dirPath", dirPath,
				"recover", err)

			// TODO: consider restarting from here, in case of crash
			// go GoMonitor(dirPath, ch)
		}
		slog.Infow("GoMonitor leaving ...", "error", err)
	}()

	// Create new watcher.
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		slog.Errorw("failed fsnotify NewWatcher", "error", err)
		return
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					slog.Errorw("failed watcher event channel closed")
					return
				}

				slog.Infow("monitor", "event", event)

				switch event.Op {
				case fsnotify.Create:
					ch <- FileInfo{PathFileName: event.Name, IsCreated: true}
				case fsnotify.Write:
					ch <- FileInfo{PathFileName: event.Name, IsModified: true}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					slog.Errorw("failed watcher errors channel closed")
					return
				}
				slog.Warnw("failed watcher errors", "error", err)
				// consider returning from here
			}
		}
	}()

	// Add a path.
	err = watcher.Add(dirPath)
	if err != nil {
		slog.Errorw("failed watcher add", "error", err)
		return
	}

	// Block this goroutine forever.
	<-make(chan struct{})
}
