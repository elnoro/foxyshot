package app

import (
	"context"
	"errors"
	"foxyshot/config"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"testing"
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
	err := testApp.Watch(ctx, ".")

	assert.Nil(t, err)
}

func initTestFoxyshotApp() *foxyshotApp {
	withS3 := &config.Config{S3: &config.S3Config{}}
	app, _ := New(withS3)

	return app.(*foxyshotApp)
}

func TestFoxyshotApp_handleEvent(t *testing.T) {
	tests := []struct {
		name        string
		ev          fsnotify.Event
		wantHandled bool
	}{
		{
			name: "create event for actual screenshot",
			ev: fsnotify.Event{
				Name: "path/to/valid-file.jpg",
				Op:   fsnotify.Create,
			},
			wantHandled: true,
		},
		{
			name: "create event for temporary screenshot file",
			ev: fsnotify.Event{
				Name: "path/to/.file-with-dot.jpg",
				Op:   fsnotify.Create,
			},
			wantHandled: false,
		},
		{
			name: "remove event",
			ev: fsnotify.Event{
				Name: "path/to/valid-file.jpg",
				Op:   fsnotify.Remove,
			},
			wantHandled: false,
		},
		{
			name: "rename event",
			ev: fsnotify.Event{
				Name: "path/to/valid-file.jpg",
				Op:   fsnotify.Rename,
			},
			wantHandled: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pipelineMock{}
			fa := &foxyshotApp{
				uploader: &uploaderMock{},
				pipeline: p,
			}
			fa.handleEvent(context.Background(), tt.ev)

			handled := p.pathCalled != ""
			assert.Equal(t, tt.wantHandled, handled)
		})
	}
}

type pipelineMock struct {
	pathCalled string
}

func (p *pipelineMock) Run(path string) (string, error) {
	p.pathCalled = path

	return path + "-processed", nil
}

type uploaderMock struct {
	pathUploaded string
}

func (u *uploaderMock) Upload(_ context.Context, path string) (string, error) {
	u.pathUploaded = path

	return "", errors.New("stop running")
}
