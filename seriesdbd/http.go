package main

import (
	"github.com/tadasv/seriesdb/util"
	"log"
	"net"
	"net/http"
)

func httpServer(listener net.Listener) {
	log.Printf("HTTP: listening on %s", listener.Addr().String())

	handler := http.NewServeMux()
	handler.HandleFunc("/alive", aliveHandler)

	server := &http.Server{
		Handler: handler,
	}
	err := server.Serve(listener)
	if err != nil {
		log.Printf("ERROR: http.Serve() - %s", err.Error())
	}

	log.Printf("HTTP: closing %s", listener.Addr().String())
}

func aliveHandler(w http.ResponseWriter, req *http.Request) {
	util.ApiResponse(w, 200, "OK", nil)
}
