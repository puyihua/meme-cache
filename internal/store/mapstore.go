package store

import (
	"fmt"
)

type MapStore struct {
	store map[string]string
}

func NewMapStore() *MapStore {
	return &MapStore{
		store: make(map[string]string),
	}
}

func (ms *MapStore) Get(key string) (string, error) {
	if val, ok := ms.store[key]; ok {
		return val, nil
	} else {
		return "", fmt.Errorf("key not found")
	}
}

func (ms *MapStore) Put(key string, value string) {
	ms.store[key] = value
}

func (ms *MapStore) GetLength() int {
	return len(ms.store)
}
