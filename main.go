package main

import (
	"log"
	"task-axxon/config"
	"task-axxon/transport"
)

func main() {
	conf := config.Load()
	server := transport.InitServer(conf)
	log.Fatal(server.ListenAndServe())
}
