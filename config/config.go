package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

// Config Main config for the application
type Config struct {
	// Folder where screenshots are stored
	WatchFor string
	S3       *S3Config
	// Compression level for JPEGs
	JpegQuality int
	// Remove original screenshot files to save space
	RemoveOriginals bool
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
	Duration int
	CDN      string
}

const (
	defaultJpegQuality = 30
	defaultBucket      = "foxy"
	defaultDuration    = 24
)

func setupViper(v *viper.Viper) {
	v.SetDefault("screenshots.jpegQuality", defaultJpegQuality)
	v.SetDefault("screenshots.removeOriginals", false)
	v.SetDefault("s3.publicURIs", true)
	v.SetDefault("s3.bucket", defaultBucket)
	v.SetDefault("s3.duration", defaultDuration)

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
	s3config := &S3Config{
		Key:        v.GetString("s3.key"),
		Secret:     v.GetString("s3.secret"),
		Endpoint:   v.GetString("s3.endpoint"),
		Region:     v.GetString("s3.region"),
		Bucket:     v.GetString("s3.bucket"),
		PublicURIs: v.GetBool("s3.publicURIs"),
		Duration:   v.GetInt("s3.duration"),
		CDN:        v.GetString("s3.cdn"),
	}

	watchFolder := expandHomeFolder(v.GetString("watchFolder"))
	config := &Config{
		WatchFor:        watchFolder,
		S3:              s3config,
		JpegQuality:     v.GetInt("screenshots.jpegQuality"),
		RemoveOriginals: v.GetBool("screenshots.removeOriginals"),
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
