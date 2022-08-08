package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDefaultSettings(t *testing.T) {
	v := viper.New()
	setupViper(v)

	assert.Equal(t, defaultJpegQuality, v.GetInt("screenshots.jpegQuality"))
	assert.Equal(t, false, v.GetBool("screenshots.removeOriginals"))
}

func TestValidConfig(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("./testdata/full.json")

	c, err := parseConfigToStruct(v)

	assert.NoError(t, err)
	assert.Equal(t, 999, c.Screenshots.JpegQuality)
	assert.Equal(t, true, c.Screenshots.RemoveOriginals)
	assert.Equal(t, "expected_folder", c.WatchFor)
	assert.Equal(t, "expected_key", c.S3.Key)
	assert.Equal(t, "expected_secret", c.S3.Secret)
	assert.Equal(t, "expected_endpoint", c.S3.Endpoint)
	assert.Equal(t, "expected_region", c.S3.Region)
	assert.Equal(t, "expected_bucket", c.S3.Bucket)
	assert.Equal(t, "expected_cdn", c.S3.CDN)
	assert.Equal(t, false, c.S3.PublicURIs)
	assert.Equal(t, 111, c.S3.Duration)
}

func TestExpandHomeFolder(t *testing.T) {
	home, _ := os.UserHomeDir()
	v := viper.New()
	v.SetConfigFile("./testdata/expand.json")
	c, err := parseConfigToStruct(v)

	assert.NoError(t, err)
	assert.Equal(t, home+"/watchfolder", c.WatchFor)
}
