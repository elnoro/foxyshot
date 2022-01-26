package cmd

import (
	"context"
	"fmt"
	"foxyshot/app"
	"foxyshot/config"
)

func startApp() error {
	appConfig := config.Load()
	cmdApp, err := app.New(appConfig)
	if err != nil {
		return fmt.Errorf("cannot start daemon, %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	go cmdApp.Watch(ctx, appConfig.WatchFor)

	cmdApp.WaitForExit(cancel)
	return nil
}
