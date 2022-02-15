package store

import (
	"fmt"
	"sync"
	"time"
)

// InMemoryData represens data in memory store
type InMemoryData struct {
	creationTime time.Time
	data         []byte
}

// InMemoryStore uses in memory
type InMemoryStore struct {
	data       map[string]InMemoryData
	expiration time.Duration
	mutex      *sync.RWMutex
}

// NewInMemoryStore create a new InMemoryStore
func NewInMemoryStore(expiration time.Duration) *InMemoryStore {
	store := &InMemoryStore{
		data:       map[string]InMemoryData{},
		expiration: expiration,
		mutex:      &sync.RWMutex{},
	}

	go func() {
		for {
			time.Sleep(store.expiration)

			keysToDelete := []string{}
			store.mutex.RLock()
			for key, data := range store.data {
				if time.Since(data.creationTime) > store.expiration {
					keysToDelete = append(keysToDelete, key)
				}
			}
			store.mutex.RUnlock()

			store.mutex.Lock()
			for _, key := range keysToDelete {
				delete(store.data, key)
			}
			store.mutex.Unlock()
		}
	}()

	return store
}

// Get data from store with given key
func (store InMemoryStore) Get(key string) ([]byte, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	if data, ok := store.data[key]; ok && time.Since(data.creationTime) <= store.expiration {
		return data.data, nil
	}

	return nil, fmt.Errorf("no data for key=%s", key)
}

// Put data to store fore given key
func (store *InMemoryStore) Set(key string, data []byte) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.data[key] = InMemoryData{
		creationTime: time.Now(),
		data:         data,
	}

	return nil
}
