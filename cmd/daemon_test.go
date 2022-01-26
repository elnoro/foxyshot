package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDaemonLifecycle(t *testing.T) {
	stateFile = "./testingdata/lifecycle.state"

	// not running yet
	p, err := getPID()
	assert.Equal(t, 0, p)
	assert.Error(t, err)
	assert.False(t, assert.FileExists(new(testing.T), stateFile))

	err = startDaemon("echo")
	assert.NoError(t, err)
	// running
	p, err = getPID()
	assert.NotEqual(t, 0, p)
	assert.NoError(t, err)
	assert.FileExists(t, stateFile)

	err = stopDaemon()
	assert.NoError(t, err)
	// no longer running
	p, err = getPID()
	assert.Equal(t, 0, p)
	assert.Error(t, err)
	assert.False(t, assert.FileExists(new(testing.T), stateFile))
}

func TestCannotStartDaemonTwice(t *testing.T) {
	dir, err := os.MkdirTemp("", "testdaemon")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(dir)

	stateFile = path.Join(dir, "test.state")

	err = startDaemon("echo")
	assert.NoError(t, err)

	err = startDaemon("echo")
	assert.Error(t, err)
	assert.True(
		t,
		strings.HasPrefix(err.Error(), "Daemon is already running"),
		fmt.Sprintf("unexpected error message: %s", err.Error()),
	)
}

func TestStartDaemon_InaccessibleLocation(t *testing.T) {
	stateFile = "./testingdata/doesnotexists/cannot.state"
	err := startDaemon("echo")
	assert.Error(t, err)

	p, err := getPID()
	assert.Equal(t, 0, p)
	assert.Error(t, err)
}

func TestGetPid_MalformedStatus(t *testing.T) {
	stateFile = "./testingdata/malformed.state"

	p, err := getPID()
	assert.Equal(t, 0, p)
	assert.Error(t, err)
}

func TestStopDaemon_DoesNotRemoveMalformedState(t *testing.T) {
	stateFile = "./testingdata/malformed.state"

	err := stopDaemon()
	assert.NoError(t, err)
	assert.FileExists(t, stateFile)
}

func TestPrintStatus_DoesNotFail(t *testing.T) {
	stateFile = "./testingdata/malformed.state"
	err := printStatus()
	assert.NoError(t, err)
	stateFile = "./testingdata/valid.state"
	err = printStatus()
	assert.NoError(t, err)

}
