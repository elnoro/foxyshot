package cmd

import (
	"context"
	"foxyshot/app"
	"foxyshot/config"
	"log"
)

func startApp() {
	appConfig := config.Load()
	app, err := app.New(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go app.Watch(ctx, appConfig.WatchFor)

	app.WaitForExit(cancel)
}
