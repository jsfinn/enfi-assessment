package mock

import (
	"errors"
	"log"
	"math/rand/v2"
	"slices"
	"strconv"
	"time"

	"github.com/jsfinn/enfi-assessment/model"
)

type fileProvider struct {
	files        []*file
	fileById     map[model.FileId]*file
	childrenById map[model.FileId][]model.FileId
}

type file struct {
	FileId       model.FileId
	LastModified int64
	IsDirectory  bool
	ParentId     model.FileId
}

// Helper function to extract metadata from a mock file
func MetadataFromFile(file file) model.Metadata {
	return model.Metadata{
		Id:           file.FileId,
		LastModified: file.LastModified,
		IsDirectory:  file.IsDirectory,
	}
}

// NewFileProvider creates a new file provider with the given number of files and directories.  The
// tree structure is randomly generated.  If a structure is needed, initialize the with 0 files and 0 directories and
// manually add the files and directories using AddFile and AddDirectory.
func NewFileProvider(fileCount int, directoryCount int) *fileProvider {

	fp := &fileProvider{fileById: make(map[model.FileId]*file), childrenById: make(map[model.FileId][]model.FileId)}

	for i := 0; i < directoryCount; i++ {
		fileId := model.FileId("directory" + strconv.Itoa(i+1))
		fp.AddDirectory(fileId, model.FileId(""))
	}

	for i := 0; i < fileCount; i++ {
		fileId := model.FileId("file" + strconv.Itoa(i+1))
		fp.AddFile(fileId, model.FileId(""))
	}

	// Randomly create a tree structure
	for i := 2; i < directoryCount; i++ {
		parentIndex := randRange(0, i)
		fp.files[i].ParentId = fp.files[parentIndex].FileId
	}

	// randomize the parent of the files
	for i := directoryCount + fileCount/10; i < fileCount+directoryCount; i++ {
		parentIndex := randRange(0, directoryCount)
		fp.files[i].ParentId = fp.files[parentIndex].FileId
	}

	return fp
}

// UpdateLastModified updates the last modified time of the file with the given ID.
func (fp *fileProvider) UpdateLastModified(fileId model.FileId) {
	if file, ok := fp.fileById[fileId]; ok {
		file.LastModified = time.Now().UnixMilli()
	}
}

func (fp *fileProvider) UpdateAny() model.FileId {
	fileIndex := randRange(0, len(fp.files))
	file := fp.files[fileIndex]
	file.LastModified = time.Now().UnixMilli()
	return file.FileId
}

// AddFile adds a file to the file provider with the given ID and parent directory.
func (fp *fileProvider) AddFile(id model.FileId, parentDirectory model.FileId) {
	millis := time.Now().UnixMilli()
	file := &file{FileId: id, LastModified: millis, IsDirectory: false, ParentId: parentDirectory}
	fp.files = append(fp.files, file)
	fp.fileById[file.FileId] = file
	fp.childrenById[file.ParentId] = append(fp.childrenById[file.ParentId], file.FileId)
}

// AddDirectory adds a directory to the file provider with the given ID and parent directory.
func (fp *fileProvider) AddDirectory(id model.FileId, parentDirectory model.FileId) {
	millis := time.Now().UnixMilli()
	directory := &file{FileId: id, LastModified: millis, IsDirectory: true, ParentId: parentDirectory}
	fp.files = append(fp.files, directory)
	fp.fileById[directory.FileId] = directory
	fp.childrenById[id] = []model.FileId{}
	fp.childrenById[directory.ParentId] = append(fp.childrenById[directory.ParentId], directory.FileId)
}

// CreateWatchList creates a watch list of the given size.  The watch list is a list of file IDs that are randomly selected from the files in the file provider.
func (fp *fileProvider) CreateWatchList(count int) []model.FileId {
	var watchList []model.FileId
	availableFiles := slices.Clone(fp.files)

	for i := 0; i < count; i++ {
		// Randomly select a file that is not in the watch list
		fileIndex := randRange(0, len(availableFiles))
		fileId := availableFiles[fileIndex].FileId
		availableFiles[i] = availableFiles[len(availableFiles)-1]
		availableFiles = availableFiles[:len(availableFiles)-1]

		watchList = append(watchList, fileId)
	}
	return watchList
}

////////////////////////////////////////
// Mock implementation of the file provider interface
////////////////////////////////////////

// RetrieveMetadata returns the metadata for the file with the given ID. If the file does not exist, it returns an error.
func (fp *fileProvider) RetrieveMetadata(fileId model.FileId) (model.Metadata, error) {
	if _, ok := fp.fileById[fileId]; !ok {
		return model.Metadata{}, errors.New("file not found")
	}
	return MetadataFromFile(*fp.fileById[fileId]), nil
}

// CopyFile copies the file with the given ID.  If the file is a directory, it returns an error.
func (fp *fileProvider) CopyFile(fileId model.FileId, version int) error {
	if file, ok := fp.fileById[fileId]; !ok {
		return errors.New("file not found")
	} else if file.IsDirectory {
		return errors.New("file is a directory")
	}
	log.Println("Copying file ", fileId, " version ", version)
	return nil
}

// GetChildren returns the children of the given file.  If FileID is empty, it returns the root directory.  If the file is not a directory, it returns an error.
func (fp *fileProvider) GetChildren(fileId model.FileId) ([]model.Metadata, error) {

	if fileIds, ok := fp.childrenById[fileId]; ok {
		var metadata []model.Metadata
		for _, fileId := range fileIds {
			file := fp.fileById[fileId]
			filemetadata, err := fp.RetrieveMetadata(file.FileId)
			if err != nil {
				return nil, err
			}
			metadata = append(metadata, filemetadata)
		}
		return metadata, nil
	}
	return nil, errors.New("file not found")
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
