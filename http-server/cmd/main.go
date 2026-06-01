package main

import (
	"log"
	"net/http"
	httpserver "test/http-server"
)

func main() {
	// handler := http.HandlerFunc(httpserver.PlayerServer)
	server := &httpserver.PlayerServer{}
	log.Fatal(http.ListenAndServe(":5000", server))
}
