package config

import (
	"log"
	"os"
	"strings"

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

func setupViper(v *viper.Viper) {
	v.SetDefault("jpegQuality", 30)

	v.SetConfigName("config")
	v.AddConfigPath("$HOME/.config/foxyshot")
	v.AddConfigPath(".")
}

func parseConfigToStruct(v *viper.Viper) *Config {
	err := v.ReadInConfig()
	if err != nil {
		// TODO ask for credentials and generate config automatically
		log.Fatalf("Cannot find the config file, got error %v", err)
	}
	creds := &S3Config{
		Key:      v.GetString("creds.key"),
		Secret:   v.GetString("creds.secret"),
		Endpoint: v.GetString("creds.endpoint"),
		Region:   v.GetString("creds.region"),
	}

	config := &Config{
		WatchFor:    v.GetString("watchFolder"),
		S3:          creds,
		JpegQuality: v.GetInt("jpegQuality"),
	}

	log.Printf("Loaded config from %s \n", viper.ConfigFileUsed())
	log.Printf("Watching folder %s. Screenshots will be uploaded to %s \n", config.WatchFor, config.S3.Endpoint)

	return config
}

func expandHomeFolder(orig string) string {
	if strings.Contains(orig, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			return strings.Replace(orig, "~", home, 1)
		}

		return orig
	}

	return orig
}

// Load Looks for config in the filesystem
func Load() *Config {
	v := viper.New()
	setupViper(v)

	return parseConfigToStruct(v)
}
