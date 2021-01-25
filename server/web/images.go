package web

import (
	"log"
	"net/http"
)

// Start starts web server to serve images to the users
func Start(address, path string) {
	log.Println("Starting the web server:", address, path)

	fs := http.FileServer(http.Dir(path))
	// TODO disable directory listing
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(address, fs))
}
