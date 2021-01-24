package web

import (
	"log"
	"net/http"
)

// StartWebServer starts web server to serve images to the users
func StartWebServer(address, path string) {
	log.Println("Starting the web server:", address, path)

	fs := http.FileServer(http.Dir(path))
	// TODO disable directory listing
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(address, fs))
}
