package main

import (
	"errors"
	"log"

	"example.com/config"
	"example.com/monitor"
	"example.com/store"

	"go.uber.org/zap"
)

func main() {
	// Note: default logging into console but complex logging to kafka into file can also be achieved using advance options
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	slog := logger.Sugar()
	defer logger.Sync()

	var conf *config.Configs
	fileChan := make(chan monitor.FileInfo, 100)

	slog.Infow("main starting ...")

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
			slog.Errorw("main application crashed", "recover", err)
		}
		slog.Infow("main leaving ...")
	}()

	conf, err = config.Config(slog)
	if err != nil {
		slog.Errorw("main configs failed", "error", err)
		return
	}

	db, err := store.Open(conf.Database.DBName, slog)
	if err != nil {
		slog.Errorw("failed store",
			"DBName", conf.Database.DBName,
			"error", err)
		return
	}

	go monitor.GoMonitor(conf.Monitor.DirPath, fileChan, slog)

	db.Save(fileChan, slog)

	// Note: Sample json log using slog.Infow():
	// {"level":"info","ts":1713825338.978036,"caller":"src/main.go:52","msg":"msg","pathFileName":"/etc/temp/newfile2","Size":71}

	// NOTE: The advantage of zap json logs (Infow, Warnw, Errorw, *w) is that they need not go through ELK stack
	// the pain sticking process of using unstructured logs into elastic search, log stash and kibana
	// instead write logs in proper json format right from begenning and easy to extract meaning out of them easily later and avoid using ELK stack.
}
