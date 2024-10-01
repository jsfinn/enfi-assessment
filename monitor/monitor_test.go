package monitor

import (
	"log"
	"testing"
	"time"

	"github.com/jsfinn/enfi-assessment/mock"
	"github.com/jsfinn/enfi-assessment/model"
)

func TestMonitor(t *testing.T) {

	fp := mock.NewFileProvider(0, 0)
	fp.AddDirectory("dir1", "")
	fp.AddDirectory("dir2", "dir1")
	fp.AddDirectory("dir3", "dir1")
	fp.AddFile("file1", "")
	fp.AddFile("file2", "dir1")
	fp.AddFile("file3", "dir2")

	historyCache := NewHistoryCache()
	watchList := []model.FileId{"dir1"}

	log.Printf("watchList: %v", watchList)

	monitor := NewMonitor(fp, watchList, historyCache)

	monitor.ExecuteTask()
	time.Sleep(10 * time.Millisecond)
	log.Printf("-------------------")
	time.Sleep(10 * time.Millisecond)
	monitor.ExecuteTask()
	log.Printf("-------------------")
	time.Sleep(10 * time.Millisecond)
	fp.UpdateLastModified("file2")
	monitor.ExecuteTask()
	log.Printf("-------------------")
	time.Sleep(10 * time.Millisecond)
	monitor.ExecuteTask()
}
