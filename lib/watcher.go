package configclient

import (
	storage "github.com/cmiceli/configserver/lib"
	"time"
)

type FileWatcher struct {
	client    storage.Storage
	files     []string
	quitChans []chan bool
}

func (f *FileWatcher) AddFile(id string, path string, t time.Duration) {
	f.files = append(f.files, path)
	f.quitChans = append(f.quitChans, make(chan bool))
	ticker := time.NewTicker(t)
	go f.watch(id, path, ticker, f.quitChans[len(f.quitChans)-1])
}

func (f *FileWatcher) watch(id string, path string, ticker *time.Ticker, quit chan bool) {
	for {
		select {
		case <-ticker.C:
			c, err := f.client.Get(id)
			if err != nil {
				panic(err)
			}
			WriteConfig(path, c)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (f *FileWatcher) Stop() {
	for _, i := range f.quitChans {
		close(i)
	}
}

func NewFileWatcher(client storage.Storage) *FileWatcher {
	return &FileWatcher{client: client}
}
