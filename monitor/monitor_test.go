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
	simpleCounter := NewSimpleCounter()

	log.Printf("watchList: %v", watchList)

	monitor := NewMonitor(fp, watchList, historyCache, simpleCounter)
	monitor.Start()

	monitor.EvaluateWatchlist()
	time.Sleep(10 * time.Millisecond)
	log.Printf("-------------------")
	time.Sleep(10 * time.Millisecond)
	monitor.EvaluateWatchlist()
	log.Printf("-------------------")
	time.Sleep(10 * time.Millisecond)
	fp.UpdateLastModified("file2")
	monitor.EvaluateWatchlist()
	log.Printf("-------------------")
	time.Sleep(10 * time.Millisecond)
	fp.UpdateLastModified("file2")
	monitor.EvaluateWatchlist()

	monitor.ShutDown()
}

func TestMonitorWithScale(t *testing.T) {
	// 1000 files, 50 directories
	fileCount := 5000
	directoryCount := 100
	watchCount := 100

	simpleCounter := NewSimpleCounter()
	fp := mock.NewFileProvider(fileCount, directoryCount)

	historyCache := NewHistoryCache()

	// 500 files to watch
	watchList := fp.CreateWatchList(watchCount)

	monitor := NewMonitor(fp, watchList, historyCache, simpleCounter)
	monitor.Start()

	// check the watchlist 100 times
	for i := 0; i < 10; i++ {
		// randomly update 100 files
		ids := []model.FileId{}
		for j := 0; j < 10; j++ {
			_ = append(ids, fp.UpdateAny())
		}
		monitor.EvaluateWatchlist()
		time.Sleep(1000 * time.Millisecond)
	}
}
