package main

import (
	"log"
	"time"

	"github.com/jsfinn/enfi-assessment/mock"
	"github.com/jsfinn/enfi-assessment/monitor"
	"github.com/spf13/viper"
)

type Config struct {
	Datafile        string `mapstructure:"datafile"`
	WatchIntervalMs int64  `mapstructure:"watch_interval_ms"`
}

func loadConfig() (*Config, error) {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // path to look for the config file in

	if err := viper.ReadInConfig(); err != nil { // Find and read the config file
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil { // Unmarshal the config into the struct
		return nil, err
	}

	return &config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fp, watchlist, steps, err := mock.NewFileProviderFromFile(config.Datafile)

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	historyCache := monitor.NewHistoryCache()

	// 500 files to watch

	monitor := monitor.NewMonitor(fp, watchlist, historyCache)
	monitor.Start()

	// check the watchlist 100 times
	for _, step := range steps {
		for _, fileId := range step {
			fp.UpdateLastModified(fileId)
		}
		monitor.EvaluateWatchlist()
		time.Sleep(time.Duration(config.WatchIntervalMs) * time.Millisecond)
	}
}
