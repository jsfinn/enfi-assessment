package model

type FileId string

type Metadata struct {
	Id           FileId
	LastModified int64
	IsDirectory  bool
}
