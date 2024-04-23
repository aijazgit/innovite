package config

import (
	"strings"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

const (
	_CONFIG_FILE_NAME_ = "config"
	_CONFIG_FILE_PATH_ = "."
	_CONFIG_FILE_EXT_  = "yaml"
)

func Config(slog *zap.SugaredLogger) (*Configs, error) {
	// Set the file name of the configurations file
	viper.SetConfigName(_CONFIG_FILE_NAME_)

	// Set the path to look for the configurations file
	viper.AddConfigPath(_CONFIG_FILE_PATH_)

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType(_CONFIG_FILE_EXT_)

	var configs Configs

	if err := viper.ReadInConfig(); err != nil {
		slog.Errorw("failed reading config", "error", err)
		return nil, err
	}

	// Set undefined variables, if required
	viper.SetDefault("database.dbname", "gorm.db")

	err := viper.Unmarshal(&configs)
	if err != nil {
		slog.Errorw("failed unmarshal config", "error", err)
		return nil, err
	}

	slog.Infow("config info", "dbname", configs.Database.DBName)
	slog.Infow("config info", "monitor dirpath", configs.Monitor.DirPath)

	// TODO: Implement viper.OnConfigChange() and viper.WatchConfig() to continiously/dynamically monitor the configuration changes
	/*
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
		})
		viper.WatchConfig()
	*/

	return &configs, nil
}
