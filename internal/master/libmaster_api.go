package master

import "hash/fnv"

// LibMaster defines the methods a master server will use
// The client can choose:
// 	1. Use the master to router, get the host owning key, and request the host
//	2. Use the master as a proxy to forward the request
type LibMaster interface {
	Router(key string) (string, error)
	Get(key string) (string, error)
	Put(key string, value string) error
	Delete(key string) error
	AddMember(hostport string, vids []uint64) error
	RemoveMember(hostport string) error
	GetMembers() []string
	Migrate(low uint64, high uint64, target string, source string) error
}

// hashKey has a given string to an unsigned 64-bit integer
func hashKey(key string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	return h.Sum64()
}

