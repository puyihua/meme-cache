package main

import (
	"net/http"
	"time"

	"github.com/puyihua/meme-cache/internal/master"
	"github.com/puyihua/meme-cache/internal/node"
)

func main() {
	// Hello world, the web server
	go func() {
		masterSrv := master.NewServerMaster(8081)
		masterSrv.Serve()
	}()

	time.Sleep(1 * time.Second)

	// register cache server on the master
	// use http.Head instead of http.Get
	// https://stackoverflow.com/questions/18598780/is-resp-body-close-necessary-if-we-dont-read-anything-from-the-body
	http.Head("http://localhost:8081/addMember?host=127.0.0.1&port=8082&vid=3885454534235")

	cacheSrv := node.NewServer(8082)
	cacheSrv.Serve()
}
