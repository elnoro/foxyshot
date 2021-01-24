package main

import (
	"context"
	"flag"
	"fmt"
	"foxyshot/server/grpcapi"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func main() {
	webAddress := flag.String("web-address", ":3322", "port for serving images")
	apiAddress := flag.String("api-address", ":3311", "port for uploading images")
	dir := flag.String("d", "/tmp/foxyshot-server", "folder to upload images to in")
	url := flag.String("u", "http://localhost:3322", "url from which images will be server")

	go startWebServer(*webAddress, *dir)
	go startAPIServer(*apiAddress, *dir, *url)

	// TODO add signal processing
	select {}
}

// startWebServer starts web server to serve images to the users
func startWebServer(address, path string) {
	log.Println("Starting the web server:", address, path)

	fs := http.FileServer(http.Dir(path))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(address, fs))

}

// startAPIServer starts GRPC server to upload images
func startAPIServer(address, path, httpURL string) {
	log.Println("Starting the api server:", address, path, httpURL)
	screenshotServer := &DiskStorage{
		path: path,
		url:  httpURL,
	}

	gs := grpc.NewServer()
	tcpListener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	grpcapi.RegisterScreenshotsServer(gs, screenshotServer)
	log.Fatal(gs.Serve(tcpListener))
}

type DiskStorage struct {
	path string
	url  string
	grpcapi.UnimplementedScreenshotsServer
}

func (ds *DiskStorage) Upload(ctx context.Context, s *grpcapi.Screenshot) (*grpcapi.ScreenshotLink, error) {
	log.Println("Received a new screenshot,", s.Extension)
	if !(s.Extension == "jpg" || s.Extension == "png") {
		return nil, fmt.Errorf("The server only accepts jpg and png")
	}

	filename := fmt.Sprintf("%s.%s", uuid.New().String(), s.Extension)
	filepath := fmt.Sprintf("%s/%s", ds.path, filename)

	if err := ioutil.WriteFile(filepath, s.Image, 0666); err != nil {
		return nil, fmt.Errorf("could now save the image: %v", err)
	}

	link := &grpcapi.ScreenshotLink{}
	link.Link = ds.url + "/" + filename

	return link, nil
}
