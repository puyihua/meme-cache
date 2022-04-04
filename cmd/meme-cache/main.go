package main

import (
	"github.com/puyihua/meme-cache/internal/master"
	"github.com/puyihua/meme-cache/internal/server"
	"net/http"
	"time"
)

func main() {
	// Hello world, the web server
	go func() {
		masterSrv := master.NewServerMaster(8081)
		masterSrv.Serve()
	} ()

	time.Sleep(1 * time.Second)

	// register cache server on the master
	http.Get("http://localhost:8081/addMember?host=127.0.0.1&port=8082&vid=3885454534235")

	cacheSrv := server.NewServer(8082)
	cacheSrv.Serve()
}
