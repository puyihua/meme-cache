package store

import "hash/fnv"

type Store interface {
	Get(key string) (string, error)
	Put(key string, value string)
	GetRange(low uint64, high uint64) map[string]string
	MigrateRecv(m map[string]string)
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


func hash2Segment(key string, segNum uint64) int {
	return int(hashKey(key) % segNum)
}

// hashKey has a given string to an unsigned 64-bit integer
func hashKey(key string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	return h.Sum64()
}
