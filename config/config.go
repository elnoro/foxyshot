package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config Main config for the application
type Config struct {
	// Folder where screenshots are stored
	WatchFor    string `mapstructure:"watchFolder"`
	S3          S3Config
	Screenshots struct {
		// Compression level for JPEGs
		JpegQuality int
		// Remove original screenshot files to save space
		RemoveOriginals bool
	}
}

// S3Config contains config for s3
// Can be used for AWS S3, Digital Ocean spaces, Google Cloud storage etc.
type S3Config struct {
	Key        string
	Secret     string
	Endpoint   string
	Region     string
	Bucket     string
	PublicURIs bool
	// Sets an expiration date for presigned url (only used is PublicURIs is set to false in s3 config)
	Duration time.Duration
	CDN      string
}

const (
	defaultJpegQuality = 30
	defaultBucket      = "foxy"
	defaultDuration    = 24 * time.Hour
)

func setupViper(v *viper.Viper) {
	v.SetDefault("screenshots.jpegQuality", defaultJpegQuality)
	v.SetDefault("screenshots.removeOriginals", true)
	v.SetDefault("s3.publicURIs", true)
	v.SetDefault("s3.bucket", defaultBucket)
	v.SetDefault("s3.duration", defaultDuration)

	v.SetConfigName("config")
	v.AddConfigPath("$HOME/.config/foxyshot")
	v.AddConfigPath(".")
}

func parseConfigToStruct(v *viper.Viper) (*Config, error) {
	err := v.ReadInConfig()
	if err != nil {
		// TODO ask for credentials and generate config automatically
		return nil, fmt.Errorf("reading config, %w", err)
	}

	var config Config
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("parsing config, %w", err)
	}
	config.WatchFor = expandHomeFolder(config.WatchFor)

	log.Printf("Loaded config from %s \n", v.ConfigFileUsed())
	log.Printf("Watching folder %s. Screenshots will be uploaded to %s \n", config.WatchFor, config.S3.Endpoint)

	return &config, nil
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
func Load() (*Config, error) {
	v := viper.New()
	setupViper(v)

	return parseConfigToStruct(v)
}
