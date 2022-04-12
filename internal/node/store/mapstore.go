package store

import (
	"fmt"
	"sync"
)

type MapStore struct {
	store map[string]string
	mutex sync.Mutex
}

func NewMapStore() *MapStore {
	return &MapStore{
		store: make(map[string]string),
	}
}

func (ms *MapStore) Get(key string) (string, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	if val, ok := ms.store[key]; ok {
		return val, nil
	} else {
		return "", fmt.Errorf("key not found")
	}
}

func (ms *MapStore) Put(key string, value string) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.store[key] = value
}

func (ms *MapStore) GetLength() int {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	return len(ms.store)
}