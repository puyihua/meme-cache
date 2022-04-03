package main

import (
	"github.com/puyihua/meme-cache/internal/master"
)

func main() {
	// Hello world, the web server
	srv := master.NewServerMaster(8081)
	srv.Serve()
}
