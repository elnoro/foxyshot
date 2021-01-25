package main

import (
	"context"
	"flag"
	"fmt"
	"foxyshot/server/db"
	"foxyshot/server/grpcapi"
	"foxyshot/server/ocr"
	"foxyshot/server/web"
	"io/ioutil"
	"log"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func main() {
	webAddress := flag.String("web-address", ":3322", "port for serving images")
	apiAddress := flag.String("api-address", ":3311", "port for uploading images")
	dir := flag.String("d", "/tmp/foxyshot-server", "folder to upload images to in")
	url := flag.String("u", "http://localhost:3322", "url from which images will be server")
	dsn := flag.String("dsn", "root:root@tcp(localhost:3306)/testscreens", "mysql to store data")
	flag.Parse()

	go web.Start(*webAddress, *dir)
	go startAPIServer(*apiAddress, *dir, *url, *dsn)

	// TODO add signal processing
	select {}
}

// startAPIServer starts GRPC server to upload images
func startAPIServer(address, path, httpURL, dsn string) {
	log.Println("Starting the api server:", address, path, httpURL)
	p, err := newProcessingPipeline(dsn)
	if err != nil {
		log.Fatal(err)
	}
	p.Start()

	screenshotServer := &DiskStorage{
		path:     path,
		url:      httpURL,
		pipeline: p,
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
	path     string
	url      string
	pipeline *processingPipeline
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
	ds.pipeline.OnImage(filepath)

	link := &grpcapi.ScreenshotLink{}
	link.Link = ds.url + "/" + filename

	return link, nil
}

type processingPipeline struct {
	ocr                  ocr.OCR
	db                   db.Manipulator
	paths                chan string
	pathsAndDescriptions chan []string
}

func newProcessingPipeline(dsn string) (*processingPipeline, error) {
	sOCR, err := ocr.NewOCR()
	if err != nil {
		return nil, err
	}
	sDB, err := db.NewSqlDb(dsn)
	if err != nil {
		return nil, err
	}

	in := make(chan string)
	out := make(chan []string)

	return &processingPipeline{
		ocr:                  sOCR,
		db:                   sDB,
		paths:                in,
		pathsAndDescriptions: out,
	}, nil
}

func (p *processingPipeline) Start() {
	go ocr.ListenOnImages(p.ocr, p.paths, p.pathsAndDescriptions)
	go db.ListenOnImages(p.db, p.pathsAndDescriptions)
}

func (p *processingPipeline) OnImage(path string) {
	p.paths <- path
}
