package store

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type WalSyncMapStore struct {
	store   map[string]string
	mutex   sync.Mutex
	logFile *os.File
}

func (ms *WalSyncMapStore) GetRange(low uint64, high uint64) map[string]string {
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

func (ms *WalSyncMapStore) MigrateRecv(m map[string]string) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	for k, v := range m {
		ms.store[k] = v
	}
}

func NewWalSyncMapStore(logFile *os.File) *WalSyncMapStore {
	return &WalSyncMapStore{
		store:   make(map[string]string),
		logFile: logFile,
	}
}

func (ms *WalSyncMapStore) Get(key string) (string, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	if val, ok := ms.store[key]; ok {
		return val, nil
	} else {
		return "", fmt.Errorf("key not found")
	}
}

func (ms *WalSyncMapStore) Put(key string, value string) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	log.Printf("put,%s,%s\n", key, value)
	ms.store[key] = value
}

func (ms *WalSyncMapStore) GetLength() int {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	return len(ms.store)
}
