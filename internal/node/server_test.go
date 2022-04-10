package node

import (
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestNodeServer(t *testing.T) {

	port := 8001
	url := "http://localhost:" + strconv.Itoa(port)

	go func() {
		svr := NewServer(port)
		svr.Serve()
	}()

	time.Sleep(1 * time.Second)

	_, err := http.Head(url + "/put?key=1&value=1")
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = http.Head(url + "/put?key=abc&value=efg")
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = http.Head(url + "/put?key=qwe&value=404")
	if err != nil {
		t.Errorf(err.Error())
	}

	resp, err := http.Get(url + "/getlen")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()

	// test get length
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(b) != "Total num of keys: 3" {
		t.Errorf("get length response incorrect: " + string(b))
	}

	// test get
	resp, err = http.Get(url + "/get?key=abc")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(b) != "efg" {
		t.Errorf("get response incorrect: " + string(b))
	}

}
