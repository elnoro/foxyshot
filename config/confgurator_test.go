package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_forceConfig(t *testing.T) {
	t.Run("creates config file if it doesn't exist", func(t *testing.T) {
		err := forceConfig("./testdata/forceConfig.json")

		assert.NoError(t, err)
		defer os.Remove("./testdata/forceConfig.json")

		contents, err := os.ReadFile("./testdata/forceConfig.json")

		assert.NoError(t, err)
		assert.Equal(t, configTemplate, string(contents))
	})

	t.Run("doesn't overwrite existing config file", func(t *testing.T) {
		// read the contents of the file before the call
		before, err := os.ReadFile("./testdata/full.json")
		assert.NotEqual(t, configTemplate, string(before))
		assert.NoError(t, err)

		err = forceConfig("./testdata/full.json")
		assert.NoError(t, err)

		after, err := os.ReadFile("./testdata/full.json")

		assert.NoError(t, err)
		assert.Equal(t, string(before), string(after))
	})
}
