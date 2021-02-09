package cmd

import (
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

	startDaemon("echo")
	// running
	p, err = getPID()
	assert.NotEqual(t, 0, p)
	assert.NoError(t, err)
	assert.FileExists(t, stateFile)

	stopDaemon()
	// no longer running
	p, err = getPID()
	assert.Equal(t, 0, p)
	assert.Error(t, err)
	assert.False(t, assert.FileExists(new(testing.T), stateFile))
}
func TestStartDaemon_InaccessibleLocation(t *testing.T) {
	stateFile = "./testingdata/doesnotexists/cannot.state"
	startDaemon("echo")

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

	stopDaemon()
	assert.FileExists(t, stateFile)
}

func TestPrintStatus_DoesNotFail(t *testing.T) {
	stateFile = "./testingdata/malformed.state"
	printStatus()
	stateFile = "./testingdata/valid.state"
	printStatus()
}
