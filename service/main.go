package main

import (
	"github.com/fatemehkarimi/link-shortener/service/api/grpc"
	_ "github.com/lib/pq"
)

func main() {
	server := grpc.NewServer()
	server.Start()
}
