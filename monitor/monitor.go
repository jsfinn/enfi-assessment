package monitor

import (
	"log"

	"github.com/jsfinn/enfi-assessment/model"
)

type Monitor struct {
	api       Api
	watchlist map[model.FileId]bool
	cache     Cache
}

// Create a new monitor with the given API and watchlist
func NewMonitor(api Api, fileIds []model.FileId, cache Cache) *Monitor {
	watchlist := make(map[model.FileId]bool)
	for _, item := range fileIds {
		watchlist[item] = true
	}

	return &Monitor{
		api:       api,
		watchlist: watchlist,
		cache:     cache,
	}
}

func (m *Monitor) ExecuteTask() {

	watchlist := []model.FileId{}
	metadatalist := []model.Metadata{}

	for key := range m.watchlist {
		watchlist = append(watchlist, key)
	}

	for len(watchlist) > 0 {
		fileId := watchlist[0]
		watchlist = watchlist[1:]

		// Retrieve the metadata for the file
		metaData, err := m.api.RetrieveMetadata(fileId)

		if err != nil {
			log.Printf("Error retrieving metadata for FileId %s: %v", fileId, err)
			continue
		}

		// If the file is a directory, add its children to the watchlist
		if metaData.IsDirectory {
			children, err := m.api.GetChildren(fileId)
			if err != nil {
				log.Printf("Error retrieving children for FileId %s: %v", fileId, err)
				continue
			}
			for _, child := range children {
				if child.IsDirectory {
					log.Printf("Adding directory %s to watchlist", child.Id)
					if _, ok := m.watchlist[child.Id]; !ok {
						watchlist = append(watchlist, child.Id)
					}
				} else {
					log.Printf("Adding file %s to metadatalist", child.Id)
					metadatalist = append(metadatalist, child)
				}
			}
		} else {
			log.Printf("Adding watched file %s to metadatalist", metaData.Id)
			metadatalist = append(metadatalist, metaData)
		}
	}

	log.Println("metadatalist: ", metadatalist)

	for _, metadata := range metadatalist {
		lastModified, version := m.cache.Get(metadata.Id)

		if metadata.LastModified > lastModified {
			log.Printf("FileId %s has been updated", metadata.Id)
			m.cache.Update(metadata.Id, metadata.LastModified)
			m.api.CopyFile(metadata.Id, version)
		} else {
			log.Printf("FileId %s has not been updated", metadata.Id)
		}
	}
}
