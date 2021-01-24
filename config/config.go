package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config Main config for the application
// TODO can probably remove this and replace it with a viper instance
type Config struct {
	// Folder where screenshots are stored
	WatchFor string
	S3       *S3Config
	Foxyshot *FoxyshotConfig
	// Compression level for JPEGs
	JpegQuality int
	// Remove original screenshot files to save space
	RemoveOriginals bool
}

// FoxyshotConfig containts config to connect to the foxyshot-server instance if you have one
// It always has priority over S3 config
type FoxyshotConfig struct {
	Address  string
	Insecure bool
}

// S3Config contains config for s3
// Can be used for AWS S3, Digital Ocean spaces, Google Cloud storage etc.
type S3Config struct {
	Key      string
	Secret   string
	Endpoint string
	Region   string
}

const (
	defaultJpegQuality = 30
)

func setupViper(v *viper.Viper) {
	v.SetDefault("screenshots.jpegQuality", defaultJpegQuality)
	v.SetDefault("screenshots.removeOriginals", false)
	v.SetDefault("foxyshot.insecure", false)

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

	s3 := &S3Config{
		Key:      v.GetString("creds.key"),
		Secret:   v.GetString("creds.secret"),
		Endpoint: v.GetString("creds.endpoint"),
		Region:   v.GetString("creds.region"),
	}

	foxyshot := &FoxyshotConfig{
		Address:  v.GetString("foxyshot.address"),
		Insecure: v.GetBool("foxyshot.insecure"),
	}

	config := &Config{
		WatchFor:        v.GetString("watchFolder"),
		S3:              s3,
		Foxyshot:        foxyshot,
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
