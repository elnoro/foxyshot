package watcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"foxyshot/config"
	"foxyshot/storage"
	"foxyshot/system/clipboard"
	"foxyshot/system/notification"

	"github.com/fsnotify/fsnotify"

	ip "foxyshot/imageprocessing"
)

// New creates dependencies and instantiates the watcher
func New(c *config.Config) (*Watcher, error) {
	uploader := storage.NewS3Uploader(&c.S3)
	pipeline := ip.NewPipeline(c)
	clipImpl := clipboard.New()
	notifier := notification.NewNotifier()

	return &Watcher{uploader: uploader, pipeline: pipeline, notifier: notifier, clipboardCopier: clipImpl}, nil
}

type notifier interface {
	Show(title, notification string) error
}

type clipboardCopier interface {
	Copy(val string) error
}

type Watcher struct {
	uploader        storage.Uploader
	pipeline        ip.ScreenshotPipeline
	clipboardCopier clipboardCopier
	notifier        notifier
}

type fileEvent struct {
	path string
}

func (fe fileEvent) Path() string {
	return fe.path
}

func (w *Watcher) onNewScreenshot(ctx context.Context, ei fileEvent) {
	log.Println("Got event:", ei)

	processed, err := w.pipeline.Run(ei.Path())
	if err != nil {
		log.Printf("Skipping %s, reason: %v\n", ei.Path(), err)

		return
	}
	url, err := w.uploader.Upload(ctx, processed)
	if err != nil {
		log.Printf("Skipping %s, reason: %v\n", ei.Path(), err)

		return
	}

	err = os.Remove(processed)
	if err != nil {
		log.Printf("Failed to remove %s, reason: %v\n", processed, err)
	}

	log.Printf("Url: %s \n", url)
	err = w.clipboardCopier.Copy(url)
	if err != nil {
		log.Printf("Could not copy the url to clipboard, got %v", err)
	}

	err = w.notifier.Show("FoxyShot", "Screenshot uploaded")
	if err != nil {
		log.Printf("Failed to display notification, got %v", err)
	}
}

func (w *Watcher) Watch(ctx context.Context, dir string) error {
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
			w.handleEvent(ctx, ev)
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

func (w *Watcher) handleEvent(ctx context.Context, event fsnotify.Event) {
	if event.Op&fsnotify.Create != fsnotify.Create {
		return
	}

	filename := filepath.Base(event.Name)
	if filename[:1] == "." {
		// this is a temporary file created by MacOS, ignore
		return
	}

	fe := fileEvent{path: event.Name}
	w.onNewScreenshot(ctx, fe)
}
