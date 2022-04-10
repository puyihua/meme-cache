package master

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
)

// LibMasterCH implements the behavior of a master with consistent hash
type LibMasterCH struct {
	hashChain    []uint64
	hash2Server  map[uint64]string
	serverHealth map[string]int // count the number of consecutive failed request
	rwLock       sync.RWMutex
}

func NewLibMasterCH() *LibMasterCH {
	return &LibMasterCH{
		hashChain:    []uint64{},
		hash2Server:  make(map[uint64]string),
		serverHealth: make(map[string]int),
	}
}

func (l *LibMasterCH) Get(key string) (string, error) {
	hostport, errRouter := l.Router(key)

	if errRouter != nil {
		return "", errRouter
	}

	resp, err := http.Get("http://" + hostport + "/get?key=" + key)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Master->Node %v: Get {%v} failed with status code %d", hostport, key, resp.StatusCode)
		return "", errors.New("get failed")
	}

	val, _ := ioutil.ReadAll(resp.Body)

	log.Printf("Get {%v, %v} from %v\n", key, string(val), hostport)
	return string(val), nil
}

func (l *LibMasterCH) Put(key string, value string) error {
	hostport, errRouter := l.Router(key)

	if errRouter != nil {
		return errRouter
	}
	key, value = url.QueryEscape(key), url.QueryEscape(value)
	resp, err := http.Get("http://" + hostport + "/put?key=" + key + "&value=" + value)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Master->Node %v: Put {%v, %v} failed with status code %d", hostport, key, value, resp.StatusCode)
		return errors.New("put failed")
	}

	log.Printf("Master: Put {%v, %v} to %v\n", key, value, hostport)
	return nil
}

func (l *LibMasterCH) Delete(key string) error {
	panic("implement me")
}

func (l *LibMasterCH) AddMember(hostport string, vids []uint64) error {
	l.rwLock.Lock()
	defer l.rwLock.Unlock()
	// The cache server is already registered
	if _, exist := l.serverHealth[hostport]; exist {
		return errors.New("server already registered")
	}

	l.serverHealth[hostport] = 0

	// generate virtual nodes
	for i := range vids {
		l.hashChain = append(l.hashChain, vids[i])
		l.hash2Server[vids[i]] = hostport
	}

	sort.Slice(l.hashChain, func(i, j int) bool { return l.hashChain[i] < l.hashChain[j] })

	// TODO: data migration

	return nil
}

func (l *LibMasterCH) RemoveMember(hostport string) error {
	panic("implement me")
}

func (l *LibMasterCH) GetMembers() []string {
	l.rwLock.RLock()
	defer l.rwLock.RUnlock()
	var members []string
	for hostport := range l.serverHealth {
		members = append(members, hostport)
	}
	return members
}

func (l *LibMasterCH) Router(key string) (string, error) {
	l.rwLock.RLock()
	defer l.rwLock.RUnlock()
	if len(l.hashChain) == 0 {
		return "", errors.New("there is no available node")
	}

	// use the binary search to the the proper vid
	id := hashKey(key)
	left, right := 0, len(l.hashChain)
	for left < right {
		mid := left + (right-left)/2
		if l.hashChain[mid] >= id {
			right = mid
		} else {
			left = mid + 1
		}
	}

	i := left
	if left >= len(l.hashChain) {
		i = 0
	}

	return l.hash2Server[l.hashChain[i]], nil
}

// generateHashKey generate the ith key of node
func generateHashKey(hostport string, i int) uint64 {
	return hashKey(hostport + "|" + strconv.Itoa(i))
}
