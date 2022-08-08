package app

import (
	"context"
	"foxyshot/config"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewApp(t *testing.T) {
	withS3 := &config.Config{S3: config.S3Config{}}
	app, err := New(withS3)

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
	withS3 := &config.Config{S3: config.S3Config{}}
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
				uploader:  &uploaderMock{},
				pipeline:  p,
				clipboard: &systemMock{},
				notifier:  &systemMock{},
			}
			fa.handleEvent(context.Background(), tt.ev)

			handled := p.pathCalled != ""
			assert.Equal(t, tt.wantHandled, handled)
		})
	}
}

func TestFoxyshotApp_onNewScreenshot_HappyPass(t *testing.T) {
	pipeline := &pipelineMock{}
	uploader := &uploaderMock{}
	system := &systemMock{}
	fa := &foxyshotApp{uploader: uploader, pipeline: pipeline, clipboard: system, notifier: system}

	fa.onNewScreenshot(context.Background(), fileEvent{path: "expected-path"})

	assert.Equal(t, "expected-path", pipeline.pathCalled)
	assert.Equal(t, "expected-path-processed", uploader.pathUploaded)
	assert.Equal(t, "expected-path-processed-uploaded", system.copiedToClipboard)
	assert.Equal(t, "Screenshot uploaded", system.notificationShown)
}

type systemMock struct {
	copiedToClipboard string
	notificationShown string
}

func (s *systemMock) Copy(val string) error {
	s.copiedToClipboard = val
	return nil
}

func (s *systemMock) Show(_, notification string) error {
	s.notificationShown = notification
	return nil
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

	return path + "-uploaded", nil
}
