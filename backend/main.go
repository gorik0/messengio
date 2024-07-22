package main

import (
	"chatapp/chat"
	"flag"
)

func main() {

	port := flag.String("port", "8080", "port to listen on")

	flag.Parse()

	hub := chat.NewHubIts()
	server := chat.NewServer(hub)
	hub.Run()

	server.Run(*port)

}
