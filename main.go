package main

import (
	"context"
	"fmt"
	"foxyshot/clipboard"
	"foxyshot/storage"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/rjeczalik/notify"
	"github.com/spf13/viper"
)

func main() {

	appConfig := loadConfig()

	uploader := storage.NewDigitalOceanUploader(
		appConfig.creds.key,
		appConfig.creds.secret,
		appConfig.creds.endpoint,
		appConfig.creds.region,
	)

	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening on events within current working directory.
	// Dispatch each create and remove events separately to c.
	if err := notify.Watch(appConfig.watchFor, c, notify.Rename); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	for {
		ei := <-c
		onNewScreenshot(ei, uploader)
	}
}

func onNewScreenshot(ei notify.EventInfo, u storage.Uploader) {
	log.Println("Got event:", ei)

	img, err := readScreenshot(ei.Path())
	if err != nil {
		log.Printf("Skipping %s, reason: %v\n", ei.Path(), err)

		return
	}
	screenshot, err := compressScreenshot(img)
	if err != nil {
		log.Printf("Could not compress image, reason: %v\n", err)
	}

	url, err := u.Upload(context.TODO(), screenshot, storage.DefaultOptions)
	os.Remove(screenshot)
	if err != nil {
		log.Printf("Skipping %s, reason: %v\n", ei.Path(), err)

		return
	}

	log.Printf("Url: %s \n", url)
	clipboard.CopyToClipboard(url)
}

func compressScreenshot(img image.Image) (string, error) {
	file, err := ioutil.TempFile("/tmp", "scst_")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Compression is just saving to jpeg with lower quality
	// With MacOS default screenshotter saves up to 90%
	log.Println("Saving compressed screenshot ", file.Name())

	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 30})
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func readScreenshot(path string) (image.Image, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

type config struct {
	watchFor string
	creds    *s3Credentials
}

type s3Credentials struct {
	key      string
	secret   string
	endpoint string
	region   string
}

func loadConfig() *config {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/config/foxyshot")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// TODO ask for credentials and generate config automatically
		log.Fatalf("Cannot find the config file, got error %v", err)
		panic(fmt.Errorf("Cannot find config file in ~/config/foxyshot/config.json, got error %v", err))
	}
	creds := &s3Credentials{
		key:      viper.GetString("creds.key"),
		secret:   viper.GetString("creds.secret"),
		endpoint: viper.GetString("creds.endpoint"),
		region:   viper.GetString("creds.region"),
	}

	config := &config{
		watchFor: viper.GetString("watchFolder"),
		creds:    creds,
	}

	log.Printf("Loaded config from %s \n", viper.ConfigFileUsed())
	log.Printf("Watching folder %s. Screenshots will be uploaded to %s \n", config.watchFor, config.creds.endpoint)

	return config
}
