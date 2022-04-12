package store

import (
	"fmt"
	"sync"
)

type storeSegment struct {
	segMutex sync.RWMutex
	segStore map[string]string
}

func newStoreSegment() *storeSegment {
	return &storeSegment{
		segStore: make(map[string]string),
	}
}

type FineGrainedMapStore struct {
	segments    []*storeSegment
	segmentNum int
}

func NewFineGrainedMapStore() *FineGrainedMapStore {
	segments := make([]*storeSegment, DefaultSegmentNumber)
	for i := 0; i < DefaultSegmentNumber; i++ {
		segments[i] = newStoreSegment()
	}
	return &FineGrainedMapStore{segments: segments, segmentNum: DefaultSegmentNumber}
}

func (ms *FineGrainedMapStore) Get(key string) (string, error) {
	segId := hash2Segment(key, uint64(ms.segmentNum))
	// only lock the segment
	ms.segments[segId].segMutex.RLock()
	defer ms.segments[segId].segMutex.RUnlock()

	val, exist := ms.segments[segId].segStore[key]

	if !exist {
		return "", fmt.Errorf("key not found")
	}

	return val, nil
}

func (ms *FineGrainedMapStore) Put(key string, value string) {
	segId := hash2Segment(key, uint64(ms.segmentNum))
	// only lock the segment
	ms.segments[segId].segMutex.Lock()
	defer ms.segments[segId].segMutex.Unlock()

	ms.segments[segId].segStore[key] = value
}

func (ms *FineGrainedMapStore) GetLength() int {
	length := 0
	for i := 0; i < ms.segmentNum; i++ {
		ms.segments[i].segMutex.RLock()
		length += len(ms.segments[i].segStore)
		ms.segments[i].segMutex.RUnlock()
	}
	return length
}
