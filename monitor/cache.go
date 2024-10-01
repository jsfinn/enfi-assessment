package monitor

import (
	"github.com/jsfinn/enfi-assessment/model"
)

// Cache is an interface that defines the methods for a cache
type Cache interface {
	// Get returns the metadata for the file with the given ID.
	Get(id model.FileId) (lastModified int64, version int)
	// Update updates the last modified time of the file with the given ID.
	Update(id model.FileId, lastModified int64) (newVersion int)
	// GetAllCacheKeys returns all the keys in the cache
	GetAllCacheKeys() []model.FileId
}

////////////////////////
// IMPLEMENTATION     //
////////////////////////

// NewHistoryCache creates a new in-memory history cache
func NewHistoryCache() *inMemoryHistoryCache {
	return &inMemoryHistoryCache{history: make(map[model.FileId]*cacheItem)}
}

type inMemoryHistoryCache struct {
	history map[model.FileId]*cacheItem
}

// History is a struct that holds the history of a file
type cacheItem struct {
	id           model.FileId
	lastModified int64
	version      int
}

func (hc *inMemoryHistoryCache) Get(id model.FileId) (lastModified int64, version int) {
	if _, ok := hc.history[id]; !ok {
		hc.history[id] = &cacheItem{id: id, lastModified: 0, version: 0}
	}
	return hc.history[id].lastModified, hc.history[id].version
}

func (hc *inMemoryHistoryCache) Update(id model.FileId, lastModified int64) (newVersion int) {
	if history, ok := hc.history[id]; ok {
		history.lastModified = lastModified
		history.version++
	} else {
		hc.history[id] = &cacheItem{id: id, lastModified: lastModified, version: 1}
	}
	return hc.history[id].version
}

func (hc *inMemoryHistoryCache) GetAllCacheKeys() []model.FileId {
	keys := make([]model.FileId, 0, len(hc.history))
	for k := range hc.history {
		keys = append(keys, k)
	}
	return keys
}
