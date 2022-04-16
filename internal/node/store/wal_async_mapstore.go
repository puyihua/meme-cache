package store

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type WalAsyncMapStore struct {
	store   map[string]string
	mutex   sync.Mutex
	logFile *os.File
	logChan chan string
}

func NewWalAsyncMapStore(logFile *os.File) *WalAsyncMapStore {
	store := &WalAsyncMapStore{
		store:   make(map[string]string),
		logFile: logFile,
		logChan: make(chan string, 1000),
	}
	go store.asyncLogWriter()
	return store
}

func (ms *WalAsyncMapStore) asyncLogWriter() {
	for record := range ms.logChan {
		log.Println(record)
	}
}

func (ms *WalAsyncMapStore) Get(key string) (string, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	if val, ok := ms.store[key]; ok {
		return val, nil
	} else {
		return "", fmt.Errorf("key not found")
	}
}

func (ms *WalAsyncMapStore) Put(key string, value string) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.logChan <- fmt.Sprintf("put,%s,%s", key, value)
	ms.store[key] = value
}

func (ms *WalAsyncMapStore) GetLength() int {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	return len(ms.store)
}
