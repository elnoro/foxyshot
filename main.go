package main

import (
	"context"
	"foxyshot/clipboard"
	"foxyshot/config"
	ip "foxyshot/imageprocessing"
	"foxyshot/storage"
	"log"
	"os"

	"github.com/rjeczalik/notify"
)

func main() {

	appConfig := config.Load()

	uploader := storage.NewS3Uploader(appConfig.S3)
	pipeline := ip.NewPipeline(appConfig)

	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening on events within current working directory.
	// Dispatch each create and remove events separately to c.
	if err := notify.Watch(appConfig.WatchFor, c, notify.Rename); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	for {
		ei := <-c
		onNewScreenshot(ei, uploader, pipeline)
	}
}

func onNewScreenshot(ei notify.EventInfo, u storage.Uploader, p ip.ScreenshotPipeline) {
	log.Println("Got event:", ei)

	processed, err := p.Run(ei.Path())
	if err != nil {
		log.Printf("Skipping %s, reason: %v\n", ei.Path(), err)

		return
	}
	url, err := u.Upload(context.Background(), processed, storage.DefaultOptions)
	os.Remove(processed)
	if err != nil {
		log.Printf("Skipping %s, reason: %v\n", ei.Path(), err)

		return
	}

	log.Printf("Url: %s \n", url)
	err = clipboard.CopyToClipboard(url)
	if err != nil {
		log.Printf("Could not copy the url to clipboard, got %v", err)
	}
}
