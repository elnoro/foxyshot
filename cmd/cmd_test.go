package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	cases := []struct {
		args         []string
		expectedMain string
		expectedSub  string
	}{
		{[]string{"noSub"}, "noSub", "status"},
		{[]string{"main", "sub"}, "main", "sub"},
		{[]string{"path with spaces", "sub"}, "path with spaces", "sub"},
		{[]string{"manyargs", "arg1", "arg2"}, "manyargs", "arg1"},
	}

	for _, c := range cases {
		t.Run(c.expectedMain, func(t *testing.T) {
			main, sub, err := parseArgs(c.args)
			assert.NoError(t, err)
			assert.NotEmpty(t, main)
			assert.Equal(t, c.expectedSub, sub)
		})
	}
}
