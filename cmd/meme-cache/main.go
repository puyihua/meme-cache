package main

import (
	"github.com/puyihua/meme-cache/internal/server"
)

func main() {
	// Hello world, the web server
	srv := server.NewServer(8080)
	srv.Serve()
}
