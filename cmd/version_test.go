package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_printVersion(t *testing.T) {
	err := printVersion()
	assert.NoError(t, err)
}
