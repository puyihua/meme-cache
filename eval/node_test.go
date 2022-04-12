package eval

import (
	"encoding/json"
	"github.com/puyihua/meme-cache/internal/node"
	"github.com/puyihua/meme-cache/internal/node/store"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

const defaultStoreType int = store.TypeBaseline
const urlStr string = "http://localhost:8082"

func TestBaselineReadThroughPut(t *testing.T) {
	evalNodeReadThroughPut(t, store.TypeBaseline)
}

func TestRWLockReadThroughPut(t *testing.T) {
	evalNodeReadThroughPut(t, store.TypeRWLock)
}

func TestFineGrainedReadThroughPut(t *testing.T) {
	evalNodeReadThroughPut(t, store.TypeFineGrained)
}

func evalNodeReadThroughPut(t *testing.T, storeType int) {
	numKeys := 200
	keys, err1 := readKeysFromJson(numKeys)
	if err1 != nil {
		t.Error(err1)
	}

	err2 := startNodeAndPutKeys(keys, storeType)
	if err2 != nil {
		t.Error(err2)
	}

	ops := []string{"GET"}
	evalNodeThroughput(t, numKeys, 1, keys, ops)
	evalNodeThroughput(t, numKeys, 5, keys, ops)
	evalNodeThroughput(t, numKeys, 10, keys, ops)
	evalNodeThroughput(t, numKeys, 20, keys, ops)
	//evalNodeThroughput(t, numKeys, 15, keys, ops)
}

func evalNodeThroughput(t *testing.T, numKeys int, numClients int, keys []string, ops []string) {

	closeChan := make(chan int)

	start := time.Now()

	for i := 0; i < numClients; i++ {
		go loadGenClient(urlStr, ops, keys, closeChan)
	}

	countSum := 0
	for i := 0; i < numClients; i++ {
		countSum += <- closeChan
	}

	expectedCount := numKeys * numClients * len(ops)

	if countSum != expectedCount {
		t.Errorf("Send %d requests, recieve %d response\n", expectedCount, countSum)
	}

	elapsed :=  int(time.Since(start) / time.Millisecond)

	t.Logf("[Read Throughput]  concurrency: %d, request: %d, elapsed: %dms, throughtput: %f \n",
		numClients, countSum, elapsed, float32(countSum) / float32(elapsed))
}

func getNumOfKeys(urlStr string) (int, error) {
	resp, err := http.Get(urlStr + "/getlen")
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	byteArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	length, err := strconv.Atoi(string(byteArr))
	if err != nil {
		return -1, err
	}
	return length, nil
}

func readKeysFromJson(numKeys int) ([]string, error) {
	// read keywords from json
	jsonFile, err := os.Open("keywords.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	var keywordsMap map[string][]string
	json.Unmarshal([]byte(bytes), &keywordsMap)
	keywords := keywordsMap["keywords"][:numKeys]

	return keywords, nil
}

func startNodeAndPutKeys(keys []string, storeType int) error {
	go func () {
		cacheNode := node.NewServerWithType(8082, storeType)
		cacheNode.Serve()
	} ()

	// wait for the server to be launched
	time.Sleep(1500 * time.Millisecond)

	// put original keys to cache
	for _, key := range keys {
		putToCache(urlStr, key, "dummy")
	}

	return nil
}

func putToCache(urlStr string, key string, value string) error {
	keyUrl := url.QueryEscape(key)
	valueUrl := url.QueryEscape(value)
	resp, err := http.Get(urlStr + "/put?key=" + keyUrl + "&value=" + valueUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Client->Master: Put failed with status code %d", resp.StatusCode)
	}

	return nil
}
