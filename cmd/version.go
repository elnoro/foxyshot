package cmd

import (
	"fmt"
	"runtime/debug"
)

var version = "development"

func displayVersion() error {
	b, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("no version info provided with the binary")
	}

	fmt.Println("Version:", version)

	for _, kv := range b.Settings {
		if kv.Key == "vcs.revision" {
			fmt.Println("Revision:", kv.Value)
		}
		if kv.Key == "vcs.time" {
			fmt.Println("Built:", kv.Value)
		}
	}

	return nil
}
