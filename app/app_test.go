package app

import (
	"errors"
	"fmt"
	"foxyshot/config"
	"github.com/fsnotify/fsnotify"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"context"
	"github.com/stretchr/testify/assert"
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

func TestFoxyshotApp_WatchCreatedFile(t *testing.T) {
	// flaky because it relies on sleep and assertEventually
	// but cannot rewrite without rewriting the app struct
	if testing.Short() {
		t.Skip()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipeline := &pipelineMock{}
	uploader := &uploaderMock{}
	testApp := &foxyshotApp{
		uploader: uploader,
		pipeline: pipeline,
	}

	dir, err := os.MkdirTemp("", "foxyshot_test")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	go testApp.Watch(ctx, dir)

	time.Sleep(1 * time.Second)
	filePath := path.Join(dir, "expected")
	_, err = os.Create(filePath)
	assert.NoError(t, err)

	assert.Eventuallyf(t, func() bool {
		return strings.Contains(pipeline.pathCalled, filePath)
	}, 1*time.Second, 500*time.Millisecond, fmt.Sprintf("pipeline must have been called, got %s", pipeline.pathCalled))
	assert.Eventuallyf(t, func() bool {
		return strings.Contains(uploader.pathUploaded, filePath+"-processed")
	}, 1*time.Second, 500*time.Millisecond, "uploader must have been called")
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
			fa.handleEvent(context.TODO(), tt.ev)

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
