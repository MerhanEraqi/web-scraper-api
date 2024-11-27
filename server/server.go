package server

import (
	"log"
	"net/http"
)

func HostStaticPages() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

	// Log server start
	log.Printf("Starting server for static files on http://localhost:8081")
    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatal(err)
    }
}