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

func (ms *MapStore) GetRange(low uint64, high uint64) map[string]string {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	m := make(map[string]string)
	for k, v := range ms.store {
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
	return m
}

func (ms *MapStore) MigrateRecv(m map[string]string) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	for k, v := range m {
		ms.store[k] = v
	}
}
