package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"foxyshot/config"
	"foxyshot/watcher"
)

func run() error {
	appConfig, err := config.Load()
	if err != nil {
		return fmt.Errorf("cannot load config, %w", err)
	}
	cmdApp, err := watcher.New(appConfig)
	if err != nil {
		return fmt.Errorf("cannot start daemon, %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := cmdApp.Watch(ctx, appConfig.WatchFor)
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	cancel()
	log.Println("Exiting...")

	return nil
}
