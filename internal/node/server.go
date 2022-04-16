package node

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/puyihua/meme-cache/internal/node/store"
)

type Server struct {
	port    int
	store   store.Store
	logFile *os.File
}

// default use of baseline store
func NewServer(port int) *Server {
	ms := store.NewMapStore()
	return &Server{port: port, store: ms}
}

// switch to different store implementation
func NewServerWithType(port int, storeImplType int) *Server {
	var ms store.Store
	switch storeImplType {
	case store.TypeBaseline:
		ms = store.NewMapStore()
	case store.TypeRWLock:
		ms = store.NewRWMapStore()
	case store.TypeFineGrained:
		ms = store.NewFineGrainedMapStore()
	case store.TypeSyncWAL:
		f, err := os.OpenFile("log_"+strconv.Itoa(port), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.SetOutput(f)
		ms = store.NewWalSyncMapStore(f)
		return &Server{port: port, store: ms, logFile: f}
	case store.TypeAsyncWAL:
		f, err := os.OpenFile("log_"+strconv.Itoa(port), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.SetOutput(f)
		ms = store.NewWalAsyncMapStore(f)
		return &Server{port: port, store: ms, logFile: f}
	}

	return &Server{port: port, store: ms}
}

func (svr *Server) getHandler(theUrl *url.URL) string {
	queryMap, err := url.ParseQuery(theUrl.RawQuery)
	if err != nil {
		return "Wrong Request Format"
	}
	if keys, ok := queryMap["key"]; !ok {
		return "Wrong Request Format"
	} else {
		value, err := svr.store.Get(keys[0])
		if err != nil {
			return "Key Not Found"
		}
		return value
	}
}

func (svr *Server) putHandler(theUrl *url.URL) string {
	queryMap, err := url.ParseQuery(theUrl.RawQuery)
	if err != nil {
		return "Wrong Request Format"
	}

	keys, ok := queryMap["key"]
	if !ok {
		return "Wrong Request Format"
	}

	values, ok := queryMap["value"]
	if !ok {
		return "Wrong Request Format"
	}
	// log.Printf("Cache server: Put {%v, %v}\n", keys[0], values[0])
	svr.store.Put(keys[0], values[0])
	return "Done"
}

func (svr *Server) migrateSendHandler(theUrl *url.URL) string {
	queryMap, err := url.ParseQuery(theUrl.RawQuery)
	if err != nil {
		return "Wrong Request Format"
	}

	lows, ok := queryMap["low"]
	if !ok {
		return "Wrong Request Format"
	}

	highs, ok := queryMap["high"]
	if !ok {
		return "Wrong Request Format"
	}

	targets, ok := queryMap["target"]
	if !ok {
		return "Wrong Request Format"
	}

	lowStr, highStr, target := lows[0], highs[0], targets[0]

	low, _ := strconv.ParseUint(lowStr, 10, 64)
	high, _ := strconv.ParseUint(highStr, 10, 64)

	m := svr.store.GetRange(low, high)

	if len(m) == 0 {
		return "Done"
	}

	payload, err := json.Marshal(m)

	if err != nil {
		return err.Error()
	}

	resp, err := http.Post("http://" + target + "/migrateRecv", "application/json",
		bytes.NewBuffer(payload))

	if err != nil {
		return err.Error()
	}

	if resp.StatusCode != 200 {
		return "migrate fails"
	}

	return "Done"
}

func (svr *Server) migrateRecvHandler(r *http.Request) string {
	var m map[string]string
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		return err.Error()
	}

	svr.store.MigrateRecv(m)

	return "Done"
}

func (svr *Server) getLenHandler() string {
	return strconv.Itoa(svr.store.GetLength())
}

func (svr *Server) Serve() {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.getHandler(r.URL))
	})

	mux.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.putHandler(r.URL))
	})

	mux.HandleFunc("/getlen", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.getLenHandler())
	})

	mux.HandleFunc("/migrate", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.migrateSendHandler(r.URL))
	})

	mux.HandleFunc("/migrateRecv", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, svr.migrateRecvHandler(r))
	})

	if svr.logFile != nil {
		defer svr.logFile.Close()
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(svr.port), mux))
}
