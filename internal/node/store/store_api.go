package store

import "hash/fnv"

type Store interface {
	Get(key string) (string, error)
	Put(key string, value string)
	GetLength() int
}

const (
	TypeBaseline = iota
	TypeRWLock
	TypeFineGrained
	TypeLockLess
	TypeSyncWAL
	TypeAsyncWAL
)

const DefaultSegmentNumber int = 8

// hashKey has a given string to an unsigned 64-bit integer
func hash2Segment(key string, segNum uint64) int {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum64() % segNum)
}
