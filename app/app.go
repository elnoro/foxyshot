package app

import (
	"context"
	"fmt"
	"foxyshot/clipboard"
	"foxyshot/config"
	"foxyshot/storage"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	ip "foxyshot/imageprocessing"
)

// App is an interface for the main app that waits for new images to appear and watches for os signals
type App interface {
	Watch(context.Context, string) error
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

type fileEvent struct {
	path string
}

func (fe fileEvent) Path() string {
	return fe.path
}

func (fa *foxyshotApp) onNewScreenshot(ctx context.Context, ei fileEvent) {
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

func (fa *foxyshotApp) Watch(ctx context.Context, dir string) error {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("cannot create watcher, %w", err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Println("got error when closing watcher, ", err)
		}
	}(watcher)
	err = watcher.Add(dir)
	if err != nil {
		return fmt.Errorf("cannot add screenshots directory, %w", err)
	}

	for {
		select {
		case ev, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			fa.handleEvent(ctx, ev)
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Println(err)
		case <-ctx.Done():
			return nil
		}
	}
}

func (fa *foxyshotApp) handleEvent(ctx context.Context, event fsnotify.Event) {
	if event.Op&fsnotify.Create != fsnotify.Create {
		return
	}

	filename := filepath.Base(event.Name)
	if filename[:1] == "." {
		// this is a temporary file created by MacOS, ignore
		return
	}

	fe := fileEvent{path: event.Name}
	fa.onNewScreenshot(ctx, fe)
}

func (fa *foxyshotApp) WaitForExit(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	cancel()
	log.Println("Exiting...")
}
