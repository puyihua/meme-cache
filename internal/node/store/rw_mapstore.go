package store

import (
	"fmt"
	"sync"
)

type RWMapStore struct{
	store map[string]string
	mutex sync.RWMutex
}

func NewRWMapStore() *RWMapStore {
	return &RWMapStore{
		store: make(map[string]string),
	}
}

func (rms *RWMapStore) Get(key string) (string, error) {
	rms.mutex.RLock()
	defer rms.mutex.RUnlock()
	if val, ok := rms.store[key]; ok {
		return val, nil
	} else {
		return "", fmt.Errorf("key not found")
	}
}

func (rms *RWMapStore) Put(key string, value string) {
	rms.mutex.Lock()
	defer rms.mutex.Unlock()
	rms.store[key] = value
}

func (rms *RWMapStore) GetLength() int {
	rms.mutex.RLock()
	defer rms.mutex.RUnlock()
	return len(rms.store)
}

