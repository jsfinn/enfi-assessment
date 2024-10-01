package monitor

import (
	"errors"
	"log"

	"github.com/jsfinn/enfi-assessment/model"
)

type Monitor struct {
	api               Api
	cache             Cache
	watchlist         map[model.FileId]bool
	evaluationChannel chan model.Metadata
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

// Start the monitor
func (m *Monitor) Start() {
	evaluationChannel := make(chan model.Metadata, 100)
	m.evaluationChannel = evaluationChannel
	go func() {
		for val := range evaluationChannel {
			m.evaluateMetadata(val)
		}
	}()
}

// Shut down the monitor and clean up resources
func (m *Monitor) ShutDown() {
	close(m.evaluationChannel)
	m.evaluationChannel = nil
}

// Evaluate the metadata for the given file.  If the file has been modified since the last evaluation,
// it will copy the file with a new verion identifier.
func (m *Monitor) evaluateMetadata(metadata model.Metadata) {
	if lastModified, _ := m.cache.Get(metadata.Id); lastModified < metadata.LastModified {
		version := m.cache.Update(metadata.Id, metadata.LastModified)
		m.api.CopyFile(metadata.Id, version)
	}
}

// Performs the main evaluation task on the watchlist.  This function will iterate over the watchlist
// and evaluate the metadata for each file.  If the file has been modified since the last evaluation, it will
// copy the file.
func (m *Monitor) EvaluateWatchlist() error {
	if m.evaluationChannel == nil {
		return errors.New("monitor not started")
	}

	// Create a local watchlist of all files and directories to evaluate
	watchlist := []model.FileId{}

	// Add all files in the configured watchlist to the local watchlist
	for key := range m.watchlist {
		watchlist = append(watchlist, key)
	}

	// iterate over the local watchlist.  Any directories found will have their children
	// directories added to the local watchlist.
	for len(watchlist) > 0 {
		fileId := watchlist[0]
		watchlist = watchlist[1:]

		// Retrieve the metadata for the file associated with the fileId
		metadata, err := m.api.RetrieveMetadata(fileId)
		if err != nil {
			log.Printf("Error retrieving metadata for FileId %s: %v", fileId, err)
			continue
		}

		// If the file is a directory
		if metadata.IsDirectory {

			// Retrieve the children of the directory
			children, err := m.api.GetChildren(fileId)
			if err != nil {
				log.Printf("Error retrieving children for FileId %s: %v", fileId, err)
				continue
			}

			for _, child := range children {
				// If the child is a directory and not already in the configured watchlist, add it to the watchlist
				if child.IsDirectory {
					if _, ok := m.watchlist[child.Id]; !ok {
						watchlist = append(watchlist, child.Id)
					}
					// If the child is a file, add it to the evaluation channel, as we've already got the metadata
				} else {
					m.evaluationChannel <- child
				}
			}
			// If the file is not a directory, add it's metadata to the evaluation channel
		} else {
			m.evaluationChannel <- metadata
		}
	}

	return nil
}
