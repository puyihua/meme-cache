package node

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Server struct {
	port  int
	store *MapStore
}

func NewServer(port int) *Server {
	ms := NewMapStore()
	return &Server{port: port, store: ms}
}

func (svr Server) getHandler(theUrl *url.URL) string {
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

func (svr Server) putHandler(theUrl *url.URL) string {
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
	svr.store.Put(keys[0], values[0])
	return "Done"
}

func (svr Server) getLenHandler() string {
	return strconv.Itoa(svr.store.GetLength())
}

func (svr Server) Serve() {
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

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(svr.port), mux))
}
