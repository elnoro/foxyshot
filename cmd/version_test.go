package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_displayVersion(t *testing.T) {
	err := printVersion()
	assert.NoError(t, err)
}
