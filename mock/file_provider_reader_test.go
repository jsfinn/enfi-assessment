package mock

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestFileProviderReader(t *testing.T) {
	filename := "../testdata.json"

	// read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	fmt.Println(data)

	var testfile testfile
	json.Unmarshal(data, &testfile)
	fmt.Println(testfile.Filesystem)
	fmt.Println(testfile.Watchlist)
	fmt.Println(testfile.Updates)
}

func TestNewFileProviderFromFile(t *testing.T) {
	filename := "../testdata.json"

	provider, watchlist, steps, err := NewFileProviderFromFile(filename)

	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	fmt.Println(provider)
	fmt.Println(watchlist)
	fmt.Println(steps)
}
