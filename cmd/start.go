package cmd

import (
	"context"
	"fmt"
	"log"

	"foxyshot/app"
	"foxyshot/config"
	"foxyshot/system/logger"
)

func startApp(args []string) error {
	logger.FromArgs(args)
	appConfig, err := config.Load()
	if err != nil {
		return fmt.Errorf("cannot load config, %w", err)
	}
	cmdApp, err := app.New(appConfig)
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

	cmdApp.WaitForExit(cancel)
	return nil
}
