package config

// Configs exported
type Configs struct {
	Database Database `yaml:"database"`
	Monitor  Monitor  `yaml:"monitor"`
}

// DatabaseConfigurations exported
type Database struct {
	DBName string `yaml:"dbname"`
}

// Monitor exported
type Monitor struct {
	DirPath string `yaml:"dirpath"`
}
