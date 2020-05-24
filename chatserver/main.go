package main

import (
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := NewServer("/entry")
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
