package store

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestStoreReadThroughput(t *testing.T) {
	// Baseline Read
	throughputBenchmark(t, []string{"GET"}, 100000, 50, TypeBaseline)
	// RW Read
	throughputBenchmark(t, []string{"GET"}, 100000, 50, TypeRWLock)
	// Fine-Grained Read
	throughputBenchmark(t, []string{"GET"}, 100000, 50, TypeFineGrained)
}

func TestStoreReadWriteThroughput(t *testing.T) {
	// Baseline Read
	throughputBenchmark(t, []string{"PUT", "GET", "PUT"}, 100000, 50, TypeBaseline)
	// RW Read
	throughputBenchmark(t, []string{"PUT", "GET", "PUT"}, 100000, 50, TypeRWLock)
	// Fine-Grained Read
	throughputBenchmark(t, []string{"PUT", "GET", "PUT"}, 100000, 50, TypeFineGrained)
}

func throughputBenchmark(t *testing.T, ops []string, numKeys int, numClients int, storeType int) {
	// generate benchmark keys
	var keys []string
	for i := 0; i < numKeys; i++ {
		keys = append(keys, strconv.Itoa(rand.Int()))
	}

	// generate random key group for different clients
	shuffledKeysGroup := make([][]string, numClients)
	shuffledKeysGroup[0] = keys
	for i := 1; i < numClients; i++ {
		shuffledKeysGroup[i] = generateShuffle(shuffledKeysGroup[i - 1])
	}

	var store Store
	switch storeType {
	case TypeBaseline:
		store = NewMapStore()
	case TypeRWLock:
		store = NewRWMapStore()
	case TypeFineGrained:
		store = NewFineGrainedMapStore()
	}

	closeChan := make(chan int)

	start := time.Now()
	for i := 0; i < numClients; i++ {
		go worker(&store, ops, shuffledKeysGroup[i], closeChan)
	}

	for i := 0; i < numClients; i++ {
		<-closeChan
	}

	count := numKeys * numClients * len(ops)
	elapsed :=  int(time.Since(start) / time.Millisecond)
	t.Logf("Throughput: %f\n", float32(count) / float32(elapsed))
}

func worker(store *Store, ops []string, keys []string, closeChan chan int) {
	for i, key := range keys {
		for _, op := range ops {
			switch op {
			case "PUT":
				(*store).Put(key, "default" + strconv.Itoa(i))
			case "GET":
				_, _ = (*store).Get(key)
			}
		}
	}
	closeChan <- 1
}

func generateShuffle(keys []string) []string {
	newKeys := make([]string, len(keys))
	copy(newKeys, keys)
	// Fisher–Yates shuffle, reference: https://yourbasic.org/golang/shuffle-slice-array/
	for i := len(newKeys) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		newKeys[i], newKeys[j] = newKeys[j], newKeys[i]
	}
	return newKeys
}



