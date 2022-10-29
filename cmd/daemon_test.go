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
	stateFile := "./testingdata/lifecycle.state"
	d := newDaemon(stateFile)

	// not running yet
	p, err := d.getPID()
	assert.Equal(t, 0, p)
	assert.Error(t, err)
	assert.False(t, assert.FileExists(new(testing.T), stateFile))

	err = d.start("echo", "hello")
	assert.NoError(t, err)
	// running
	p, err = d.getPID()
	assert.NotEqual(t, 0, p)
	assert.NoError(t, err)
	assert.FileExists(t, stateFile)

	err = d.stop()
	assert.NoError(t, err)
	// no longer running
	p, err = d.getPID()
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

	stateFile := path.Join(dir, "test.state")
	d := newDaemon(stateFile)

	err = d.start("echo", "hello")
	assert.NoError(t, err)

	err = d.start("echo", "hello")
	assert.Error(t, err)
	assert.True(
		t,
		strings.HasPrefix(err.Error(), "Daemon is already running"),
		fmt.Sprintf("unexpected error message: %s", err.Error()),
	)
}

func TestStartDaemon_InaccessibleLocation(t *testing.T) {
	stateFile := "./testingdata/doesnotexists/cannot.state"
	d := newDaemon(stateFile)
	err := d.start("echo", "hello")
	assert.Error(t, err)

	p, err := d.getPID()
	assert.Equal(t, 0, p)
	assert.Error(t, err)
}

func TestStartDaemon_InvalidCommand(t *testing.T) {
	d := newDaemon("./testingdata/invalidcommand")

	err := d.start("")
	assert.ErrorContains(t, err, "exec: no command")
}

func TestGetPid_MalformedStatus(t *testing.T) {
	stateFile := "./testingdata/malformed.state"
	d := newDaemon(stateFile)

	p, err := d.getPID()
	assert.Equal(t, 0, p)
	assert.Error(t, err)
}

func TestStopDaemon_DoesNotRemoveMalformedState(t *testing.T) {
	stateFile := "./testingdata/malformed.state"
	d := newDaemon(stateFile)

	err := d.stop()
	assert.EqualError(t, err, "Cannot find the state of the app. Got unexpected pid value 0, reason expected integer")
	assert.FileExists(t, stateFile)
}

func TestPrintStatus_DoesNotPanic(t *testing.T) {
	stateFile := "./testingdata/malformed.state"
	err := printStatus(newDaemon(stateFile))
	assert.EqualError(t, err, "Printing status error, unexpected pid value 0, reason expected integer")
	stateFile = "./testingdata/valid.state"
	err = printStatus(newDaemon(stateFile))
	assert.NoError(t, err)
}

func Test_NewDefault(t *testing.T) {
	assert.NotNil(t, newDefault())
}
