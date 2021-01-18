package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config Main config for the application
type Config struct {
	// Folder where screenshots are stored
	WatchFor string
	S3       *S3Config
	// Compression level for JPEGs
	JpegQuality int
}

// S3Config contains config for s3
// Can be used for AWS S3, Digital Ocean spaces, Google Cloud storage etc.
type S3Config struct {
	Key      string
	Secret   string
	Endpoint string
	Region   string
}

// Load Looks for config in the filesystem
func Load() *Config {
	viper.SetDefault("jpegQuality", 30)

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/foxyshot")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// TODO ask for credentials and generate config automatically
		log.Fatalf("Cannot find the config file, got error %v", err)
		panic(fmt.Errorf("Cannot find config file in ~/.config/foxyshot/config.json, got error %v", err))
	}
	creds := &S3Config{
		Key:      viper.GetString("creds.key"),
		Secret:   viper.GetString("creds.secret"),
		Endpoint: viper.GetString("creds.endpoint"),
		Region:   viper.GetString("creds.region"),
	}

	config := &Config{
		WatchFor:    viper.GetString("watchFolder"),
		S3:          creds,
		JpegQuality: viper.GetInt("jpegQuality"),
	}

	log.Printf("Loaded config from %s \n", viper.ConfigFileUsed())
	log.Printf("Watching folder %s. Screenshots will be uploaded to %s \n", config.WatchFor, config.S3.Endpoint)

	return config
}
