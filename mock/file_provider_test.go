package mock

import (
	"fmt"
	"testing"
)

// assertEqual checks if two integers are equal.
func assertEqual(t *testing.T, got, want any, message string) {
	if got != want {
		t.Errorf("%s: got %d, want %d", message, got, want)
	}
}

func TestInitialize(t *testing.T) {
	fileCount := 500
	directoryCount := 10

	fp := NewFileProvider(fileCount, directoryCount)
	files := fp.files

	assertEqual(t, fileCount+directoryCount, len(files), "file count")
}

func TestManualInitialization(t *testing.T) {
	fp := NewFileProvider(0, 0)

	fp.AddDirectory("dir1", "")
	fp.AddDirectory("dir2", "dir1")
	fp.AddDirectory("dir3", "dir1")
	fp.AddFile("file1", "")
	fp.AddFile("file2", "dir1")
	fp.AddFile("file3", "dir2")

	assertEqual(t, 6, len(fp.files), "file count")

	rootChildren, err := fp.GetChildren("")
	assertEqual(t, nil, err, "rootChildren error")
	assertEqual(t, 2, len(rootChildren), "rootChildren")

	dir1Children, err := fp.GetChildren("dir1")
	assertEqual(t, nil, err, "dir1Children error")
	assertEqual(t, 3, len(dir1Children), "dir1Children")

	dir2Children, err := fp.GetChildren("dir2")
	assertEqual(t, nil, err, "dir2Children error")
	assertEqual(t, 1, len(dir2Children), "dir2Children")

	dir3Children, err := fp.GetChildren("dir3")
	assertEqual(t, nil, err, "dir3Children error")
	assertEqual(t, 0, len(dir3Children), "dir3Children")
}

func TestCreateWatchList(t *testing.T) {
	fileCount := 20
	directoryCount := 10
	watchCount := 10

	fp := NewFileProvider(fileCount, directoryCount)
	watchList := fp.CreateWatchList(watchCount)

	assertEqual(t, watchCount, len(watchList), "watchList count")
	fmt.Printf("watchList: %v\n", watchList)
}
