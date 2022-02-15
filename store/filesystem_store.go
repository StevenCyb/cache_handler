package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// FilesystemData represens data in file on filesystem
type FilesystemData struct {
	creationTime time.Time
	path         string
}

// FileStore uses filesystem to store data
type FilesystemStore struct {
	basePath   string
	fileIndex  map[string]FilesystemData
	expiration time.Duration
	mutex      *sync.RWMutex
}

// NewFilesystem create a new FilesystemStore
func NewFilesystem(basePath string, expiration time.Duration) *FilesystemStore {
	store := &FilesystemStore{
		basePath:   basePath,
		fileIndex:  map[string]FilesystemData{},
		expiration: expiration,
		mutex:      &sync.RWMutex{},
	}

	go func() {
		for {
			time.Sleep(store.expiration)

			keysToDelete := map[string]FilesystemData{}
			store.mutex.RLock()
			for key, fileIndex := range store.fileIndex {
				if time.Since(fileIndex.creationTime) > store.expiration {
					keysToDelete[key] = fileIndex
				}
			}
			store.mutex.RUnlock()

			store.mutex.Lock()
			for key := range keysToDelete {
				delete(store.fileIndex, key)
			}
			store.mutex.Unlock()
			for _, value := range keysToDelete {
				os.Remove(value.path)
			}
		}
	}()

	return store
}

// Get data from store with given key
func (store FilesystemStore) Get(key string) ([]byte, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	if data, ok := store.fileIndex[key]; ok && time.Since(data.creationTime) <= store.expiration {
		content, err := ioutil.ReadFile(data.path)
		if err != nil {
			return nil, err
		}
		return content, nil
	}

	return nil, fmt.Errorf("no data for key=%s", key)
}

// Put data to store fore given key
func (store *FilesystemStore) Set(key string, data []byte) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	path := store.basePath + "/" + key
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	store.fileIndex[key] = FilesystemData{
		creationTime: time.Now(),
		path:         path,
	}

	return nil
}
