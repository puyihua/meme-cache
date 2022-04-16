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

func (rms *RWMapStore) GetRange(low uint64, high uint64) map[string]string {
	rms.mutex.RLock()
	m := make(map[string]string)
	for k, v := range rms.store {
		hashValue := hashKey(k)
		if low > high {
			if hashValue < high || hashValue >= low {
				m[k] = v
			}
		} else {
			if hashValue >= low && hashValue < high {
				m[k] = v
			}
		}
	}
	rms.mutex.RUnlock()

	defer func() {
		rms.mutex.Lock()
		for key := range m {
			delete(rms.store, key)
		}
		rms.mutex.Unlock()
	} ()
	return m
}

func (rms *RWMapStore) MigrateRecv(m map[string]string) {
	rms.mutex.Lock()
	defer rms.mutex.Unlock()
	for k, v := range m {
		rms.store[k] = v
	}
}

