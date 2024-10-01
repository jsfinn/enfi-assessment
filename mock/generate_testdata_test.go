package mock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

// File represents a file or directory in the filesystem
type File struct {
	FileID      string `json:"fileId"`
	IsDirectory bool   `json:"isDirectory,omitempty"`
	Children    []File `json:"children,omitempty"`
}

// Data represents the overall JSON structure
type Data struct {
	Filesystem []File     `json:"filesystem"`
	Watchlist  []string   `json:"watchlist"`
	Updates    [][]string `json:"updates"`
}

// TestGenerateSmallTestData generates a small test data file
func TestGenerateLargeTestData(t *testing.T) {

	// Configuration parameters
	numFiles := 10000
	numDirs := 100
	watchlistSize := 500
	numIterations := 10
	updateSize := 5000

	// Generate filesystem
	filesystem, allFiles := generateFilesystem(numFiles, numDirs)

	// Generate watchlist
	watchlist := generateWatchlist(allFiles, watchlistSize)

	// Generate updates
	updates := generateUpdates(allFiles, numIterations, updateSize)

	// Compile data
	data := Data{
		Filesystem: filesystem,
		Watchlist:  watchlist,
		Updates:    updates,
	}

	// Write to JSON file
	outputFile := "../testdatalarge.json"
	err := writeJSONToFile(data, outputFile)
	if err != nil {
		t.Fatalf("Error writing JSON to file: %v", err)
	}

	fmt.Printf("%s has been generated successfully.\n", outputFile)
}

// generateFilesystem creates the filesystem structure and returns the root files and all file IDs
func generateFilesystem(numFiles, numDirs int) ([]File, []string) {
	var filesystem []File
	fileCounter := 1
	dirCounter := 1
	var allFiles []string

	for i := 0; i < numDirs; i++ {
		dirID := fmt.Sprintf("dir%d", dirCounter)
		dirCounter++
		directory := File{
			FileID:      dirID,
			IsDirectory: true,
			Children:    []File{},
		}

		// Assign a random number of files to each directory (5 to 15)
		numFilesInDir := rand.Intn(11) + 5 // 5 to 15
		for j := 0; j < numFilesInDir && fileCounter <= numFiles; j++ {
			fileID := fmt.Sprintf("file%d", fileCounter)
			fileCounter++
			directory.Children = append(directory.Children, File{FileID: fileID})
			allFiles = append(allFiles, fileID)
		}

		// Optionally, add subdirectories (1 to 3)
		numSubdirs := rand.Intn(3) + 1 // 1 to 3
		for k := 0; k < numSubdirs && dirCounter <= numDirs; k++ {
			subdirID := fmt.Sprintf("dir%d", dirCounter)
			dirCounter++
			subdir := File{
				FileID:      subdirID,
				IsDirectory: true,
				Children:    []File{},
			}

			// Assign files to subdirectories (5 to 10)
			numFilesInSubdir := rand.Intn(6) + 5 // 5 to 10
			for l := 0; l < numFilesInSubdir && fileCounter <= numFiles; l++ {
				fileID := fmt.Sprintf("file%d", fileCounter)
				fileCounter++
				subdir.Children = append(subdir.Children, File{FileID: fileID})
				allFiles = append(allFiles, fileID)
			}

			directory.Children = append(directory.Children, subdir)
		}

		filesystem = append(filesystem, directory)
	}

	// Add remaining files at root level
	for fileCounter <= numFiles {
		fileID := fmt.Sprintf("file%d", fileCounter)
		filesystem = append(filesystem, File{FileID: fileID})
		allFiles = append(allFiles, fileID)
		fileCounter++
	}

	return filesystem, allFiles
}

// generateWatchlist selects a random subset of files for the watchlist
func generateWatchlist(allFiles []string, watchlistSize int) []string {
	if watchlistSize > len(allFiles) {
		watchlistSize = len(allFiles)
	}
	// Shuffle and select the first 'watchlistSize' files
	rand.Shuffle(len(allFiles), func(i, j int) {
		allFiles[i], allFiles[j] = allFiles[j], allFiles[i]
	})
	return allFiles[:watchlistSize]
}

// generateUpdates creates update iterations with random files
func generateUpdates(allFiles []string, numIterations, updateSize int) [][]string {
	var updates [][]string
	for i := 0; i < numIterations; i++ {
		// Shuffle and select 'updateSize' files
		shuffled := make([]string, len(allFiles))
		copy(shuffled, allFiles)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		currentUpdateSize := updateSize
		if updateSize > len(shuffled) {
			currentUpdateSize = len(shuffled)
		}
		updates = append(updates, shuffled[:currentUpdateSize])
	}
	return updates
}

// writeJSONToFile marshals the data to JSON and writes it to a file
func writeJSONToFile(data Data, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // For pretty-printing
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}
