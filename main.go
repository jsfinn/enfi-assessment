package main

import (
	"log"
	"time"

	"github.com/jsfinn/enfi-assessment/mock"
	"github.com/jsfinn/enfi-assessment/model"
	"github.com/jsfinn/enfi-assessment/monitor"
	"github.com/spf13/viper"
)

type Config struct {
	FileCount             int   `mapstructure:"file_count"`
	DirectoryCount        int   `mapstructure:"directory_count"`
	WatchfileSize         int   `mapstructure:"watchfile_size"`
	WatchIterations       int   `mapstructure:"watch_iterations"`
	WatchIntervalMs       int   `mapstructure:"watch_interval_ms"`
	MutationsPerIteration int64 `mapstructure:"mutations_per_iteration"`
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

	fp := mock.NewFileProvider(config.FileCount, config.DirectoryCount)

	historyCache := monitor.NewHistoryCache()

	// 500 files to watch
	watchList := fp.CreateWatchList(config.WatchfileSize)

	monitor := monitor.NewMonitor(fp, watchList, historyCache)
	monitor.Start()

	// check the watchlist 100 times
	for i := 0; i < config.WatchIterations; i++ {
		// randomly update 100 files
		ids := []model.FileId{}
		for j := 0; j < int(config.MutationsPerIteration); j++ {
			_ = append(ids, fp.UpdateAny())
		}
		monitor.EvaluateWatchlist()
		time.Sleep(time.Duration(config.WatchIntervalMs) * time.Millisecond)
	}
}
