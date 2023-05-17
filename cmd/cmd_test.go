package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	err := RunCmd([]string{"main", "unknown-command"})

	assert.ErrorIs(t, err, errUnknownSubCommand)
}

func TestGetExecutable(t *testing.T) {
	ex := getExecutable()

	assert.NotEmpty(t, ex)
}

func TestParseArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"two valid args", args{args: []string{"arg", "expected-subcommand"}}, "expected-subcommand"},
		{"one arg - expect help", args{args: []string{"arg"}}, "help"},
		{"no args - expect help", args{}, "help"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, parseArgs(tt.args.args), "parseArgs(%v)", tt.args.args)
		})
	}
}

func Test_printHelp(t *testing.T) {
	err := printHelp()

	assert.NoError(t, err)
}
