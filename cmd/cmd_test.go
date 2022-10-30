package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	sub := parseArgs([]string{"arg", "expected-subcommand"})

	assert.Equal(t, "expected-subcommand", sub)
}

func TestRunCmd(t *testing.T) {
	err := RunCmd([]string{"main", "unknown-command"})

	assert.ErrorIs(t, err, errUnknownSubCommand)
}

func TestGetExecutable(t *testing.T) {
	ex := getExecutable()

	assert.NotEmpty(t, ex)
}
