package monitor

import (
	"github.com/jsfinn/enfi-assessment/model"
)

// History is a struct that holds the history of a file
type History struct {
	id           model.FileId
	lastModified int64
	version      int
}

// Cache is an interface that defines the methods for a cache
type Cache interface {
	// Get returns the metadata for the file with the given ID.
	Get(id model.FileId) (lastModified int64, version int)
	// Update updates the last modified time of the file with the given ID.
	Update(id model.FileId, lastModified int64) (newVersion int)
}

type HistoryCache struct {
	history map[model.FileId]*History
}

func NewHistoryCache() *HistoryCache {
	return &HistoryCache{history: make(map[model.FileId]*History)}
}

func (hc *HistoryCache) Get(id model.FileId) (lastModified int64, version int) {
	if _, ok := hc.history[id]; !ok {
		hc.history[id] = &History{id: id, lastModified: 0, version: 0}
	}
	return hc.history[id].lastModified, hc.history[id].version
}

func (hc *HistoryCache) Update(id model.FileId, lastModified int64) (newVersion int) {
	if history, ok := hc.history[id]; ok {
		history.lastModified = lastModified
		history.version++
	} else {
		hc.history[id] = &History{id: id, lastModified: lastModified, version: 1}
	}
	return hc.history[id].version
}
