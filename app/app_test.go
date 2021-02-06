package app

import (
	"foxyshot/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNewApp(t *testing.T) {
	notInitialized := &config.Config{}
	app, err := New(notInitialized)

	assert.Error(t, err)
	assert.Nil(t, app)

	withS3 := &config.Config{S3: &config.S3Config{}}
	app, err = New(withS3)

	assert.NoError(t, err)
	assert.IsType(t, &foxyshotApp{}, app)
}

func TestFoxyshotApp_WatchCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	testApp := initTestFoxyshotApp()
	cancel()

	// must return immediately, since the ctx is cancelled
	testApp.Watch(ctx, ".")
}

func initTestFoxyshotApp() *foxyshotApp {
	withS3 := &config.Config{S3: &config.S3Config{}}
	app, _ := New(withS3)

	return app.(*foxyshotApp)
}
