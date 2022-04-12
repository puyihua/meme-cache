package eval

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// concurrent clients that can generate load
func loadGenClient(urlStr string, ops []string, keys []string, closeChan chan int) {
	successfulReqCount := 0
	for i, key := range keys {
		for _, op := range ops {
			switch op {
			case "PUT":
				keyUrl := url.QueryEscape(key)
				valueUrl := url.QueryEscape("dummy" + fmt.Sprint(i))
				_, err := http.Get(urlStr + "/put?key=" + keyUrl + "&value=" + valueUrl)
				if err == nil {
					successfulReqCount += 1
				}
			case "GET":
				keyUrl := url.QueryEscape(key)
				_, err := http.Get(urlStr + "/get?key=" + keyUrl)
				if err == nil {
					successfulReqCount += 1
				}
			}
			// simulate network latency
			latency := 20 + rand.Intn(30)
			time.Sleep(time.Duration(latency) * time.Millisecond)
		}
	}
	closeChan <- successfulReqCount
}

