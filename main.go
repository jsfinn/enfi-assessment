package main

import (
	"log"
	"time"

	"github.com/jsfinn/enfi-assessment/mock"
	"github.com/jsfinn/enfi-assessment/model"
	"github.com/jsfinn/enfi-assessment/monitor"
	"github.com/samber/lo"
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

	counter := monitor.NewSimpleCounter()

	// 500 files to watch

	monitor := monitor.NewMonitor(fp, watchlist, historyCache, counter)
	monitor.Start()

	// check the watchlist 100 times
	for _, step := range steps {
		for _, fileId := range step {
			fp.UpdateLastModified(fileId)
		}
		monitor.EvaluateWatchlist()
		time.Sleep(time.Duration(config.WatchIntervalMs) * time.Millisecond)
	}
	monitor.ShutDown()

	// Dump watch Log
	log.Printf("watch Log:")
	watchlistMap := lo.Associate(watchlist, func(fileId model.FileId) (model.FileId, bool) { return fileId, true })
	historyKeys := historyCache.GetAllCacheKeys()

	for _, key := range historyKeys {
		_, version := historyCache.Get(key)
		var watchtype string
		if _, ok := watchlistMap[model.FileId(key)]; ok {
			watchtype = "explicit"
		} else {
			watchtype = "implicit"
		}
		var status = "not copied"
		if version > 0 {
			status = "copied"
		}

		log.Printf("File: %v   watchtype: %v  version: %v   status: %v", key, watchtype, version, status)
		delete(watchlistMap, key)
	}

	for key := range watchlistMap {
		if m, _ := fp.RetrieveMetadata(key); !m.IsDirectory {
			log.Printf("File: %v   watchtype: explicit  version: 0   status: not copied", key)
		}
	}

	// Dump counter stats
	counter.DumpStatsToLog()

}
