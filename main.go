package main

import (
	"context"
	"foxyshot/clipboard"
	"foxyshot/config"
	ip "foxyshot/imageprocessing"
	"foxyshot/storage"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rjeczalik/notify"
)

func main() {
	appConfig := config.Load()

	uploader := storage.NewS3Uploader(appConfig.S3)
	pipeline := ip.NewPipeline(appConfig)

	app := &foxyshotApp{uploader: uploader, pipeline: pipeline}
	ctx, cancel := context.WithCancel(context.Background())
	go app.Watch(ctx, appConfig.WatchFor)

	app.WaitForExit(cancel)
}

type foxyshotApp struct {
	uploader storage.Uploader
	pipeline ip.ScreenshotPipeline
}

func (fa *foxyshotApp) onNewScreenshot(ctx context.Context, ei notify.EventInfo) {
	log.Println("Got event:", ei)

	processed, err := fa.pipeline.Run(ei.Path())
	if err != nil {
		log.Printf("Skipping %s, reason: %v\n", ei.Path(), err)

		return
	}
	url, err := fa.uploader.Upload(ctx, processed, storage.DefaultOptions)
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

func (fa *foxyshotApp) Watch(ctx context.Context, dir string) {
	eventsChannel := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening on events within current working directory.
	// Dispatch each create and remove events separately to c.
	if err := notify.Watch(dir, eventsChannel, notify.Rename); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(eventsChannel)

	for event := range eventsChannel {
		fa.onNewScreenshot(ctx, event)
	}
}

func (fa *foxyshotApp) WaitForExit(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	cancel()
	log.Println("Exiting...")
}
