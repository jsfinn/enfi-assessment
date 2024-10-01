package mock

import (
	"encoding/json"
	"os"

	"github.com/jsfinn/enfi-assessment/model"
)

type fileDescription struct {
	Id          string             `json:"fileId"`
	IsDirectory bool               `json:"isDirectory"`
	Children    []*fileDescription `json:"children"`
}

type testfile struct {
	Filesystem []*fileDescription `json:"filesystem"`
	Watchlist  []string           `json:"watchlist"`
	Updates    [][]string         `json:"updates"`
}

func NewFileProviderFromFile(filename string) (fileProvider *fileProvider, watchlist []model.FileId, updates [][]model.FileId, err error) {
	var testfile testfile
	var data []byte

	data, err = os.ReadFile(filename)
	if err != nil {
		return
	}
	json.Unmarshal(data, &testfile)

	fileProvider = NewFileProvider(0, 0)

	for _, f := range testfile.Filesystem {
		if f.IsDirectory {
			fileProvider.AddDirectory(model.FileId(f.Id), "")
			addChildren(fileProvider, f)
		} else {
			fileProvider.AddFile(model.FileId(f.Id), "")
		}
	}

	for _, f := range testfile.Watchlist {
		watchlist = append(watchlist, model.FileId(f))
	}

	for _, u := range testfile.Updates {
		var update []model.FileId
		for _, id := range u {
			update = append(update, model.FileId(id))
		}
		updates = append(updates, update)
	}

	return
}

func addChildren(fp *fileProvider, f *fileDescription) {
	parentId := model.FileId(f.Id)
	for _, c := range f.Children {
		if c.IsDirectory {
			fp.AddDirectory(model.FileId(c.Id), parentId)
			addChildren(fp, c)
		} else {
			fp.AddFile(model.FileId(c.Id), parentId)
		}
	}
}
