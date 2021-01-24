package main

import (
	"context"
	"flag"
	"fmt"
	"foxyshot/server/grpcapi"
	"io/ioutil"
	"log"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func main() {
	grpcAddress := flag.String("a", ":3311", "port to accept client connection on")
	f := flag.String("f", "/tmp/foxyshot-server", "folder to images in")
	u := flag.String("u", "http://localhost/images", "url prefix for image links")

	screenshotServer := &DiskStorage{
		path: *f,
		url:  *u,
	}

	gs := grpc.NewServer()
	tcpListener, err := net.Listen("tcp", *grpcAddress)
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
