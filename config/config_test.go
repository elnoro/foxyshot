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

	assert.NotEqual(t, 0, v.GetInt("jpegQuality"))
}

func TestValidConfig(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("./testdata/config.json")

	c := parseConfigToStruct(v)

	assert.Equal(t, 999, c.JpegQuality)
	assert.Equal(t, "expected_folder", c.WatchFor)
	assert.Equal(t, "expected_key", c.S3.Key)
	assert.Equal(t, "expected_secret", c.S3.Secret)
	assert.Equal(t, "expected_endpoint", c.S3.Endpoint)
	assert.Equal(t, "expected_region", c.S3.Region)
}

func TestExpandHomeFolder(t *testing.T) {
	home, _ := os.UserHomeDir()
	expanded := expandHomeFolder("~/watchfolder")

	assert.Equal(t, home+"/watchfolder", expanded)
}
