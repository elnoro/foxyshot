package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_displayVersion(t *testing.T) {
	err := displayVersion()
	assert.NoError(t, err)
}
