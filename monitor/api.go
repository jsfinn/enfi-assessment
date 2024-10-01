package monitor

import "github.com/jsfinn/enfi-assessment/model"

type Api interface {
	// RetrieveMetadata returns the metadata for the file with the given ID. If the file does not exist, it returns an error.
	RetrieveMetadata(fileId model.FileId) (model.Metadata, error)

	// CopyFile copies the file with the given ID.  If the file is a directory, it returns an error.
	CopyFile(fileId model.FileId, version int) error

	// GetChildren returns the children of the given file.  If FileID is empty, it returns the root directory.  If the file is not a directory, it returns an error.
	GetChildren(fileId model.FileId) ([]model.Metadata, error)
}
