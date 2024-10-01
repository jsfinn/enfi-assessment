package monitor

import (
	"log"

	"golang.org/x/sync/syncmap"
)

type SimpleCounter struct {
	stats syncmap.Map
}

func NewSimpleCounter() *SimpleCounter {
	return &SimpleCounter{
		stats: syncmap.Map{},
	}
}

func (sc *SimpleCounter) IncrementStat(name string) {
	if _, ok := sc.stats.Load(name); !ok {
		sc.stats.Store(name, 0)
	}
	value, _ := sc.stats.Load(name)
	sc.stats.Store(name, value.(int)+1)
}

func (sc *SimpleCounter) DumpStatsToLog() {
	sc.stats.Range(func(key, value interface{}) bool {
		log.Printf("%v: %v", key, value)
		return true
	})
}
