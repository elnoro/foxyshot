package app

import (
	"context"
	"fmt"
	"foxyshot/clipboard"
	"foxyshot/config"
	"foxyshot/storage"
	"log"
	"os"
	"os/signal"
	"syscall"

	ip "foxyshot/imageprocessing"

	"github.com/rjeczalik/notify"
)

// App is an interface for the main app that waits for new images to appear and watches for os signals
type App interface {
	Watch(context.Context, string)
	WaitForExit(context.CancelFunc)
}

// New creates dependencis and instantiates the app
func New(c *config.Config) (App, error) {
	if c.S3 == nil {
		return nil, fmt.Errorf("Only S3 config is supported for now, must initialized")
	}
	uploader := storage.NewS3Uploader(c.S3)
	pipeline := ip.NewPipeline(c)

	return &foxyshotApp{uploader: uploader, pipeline: pipeline}, nil
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
	url, err := fa.uploader.Upload(ctx, processed)
	if err != nil {
		log.Printf("Skipping %s, reason: %v\n", ei.Path(), err)

		return
	}

	err = os.Remove(processed)
	if err != nil {
		log.Printf("Failed to remove %s, reason: %v\n", processed, err)
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

	for {
		select {
		case event := <-eventsChannel:
			fa.onNewScreenshot(ctx, event)
		case <-ctx.Done():
			log.Printf("Watch stopped, got %v\n", ctx.Err())
			return
		}
	}
}

func (fa *foxyshotApp) WaitForExit(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	cancel()
	log.Println("Exiting...")
}
