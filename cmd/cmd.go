package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"foxyshot/config"
)

var version = "development"

// RunCmd parses the subcommand and chooses the behaviour
func RunCmd(args []string) error {
	mainCmd, subCmd, err := parseArgs(args)
	if err != nil {
		return err
	}
	switch subCmd {
	case "run":
		return startApp(os.Args[2:])
	case "start":
		return startDaemon(mainCmd)
	case "stop":
		return stopDaemon()
	case "status":
		return printStatus()
	case "configure":
		return config.RunConfigure()
	case "version":
		return getVersion()
	default:
		return fmt.Errorf("Unknown subcommand")
	}
}

func parseArgs(args []string) (string, string, error) {
	mainCmd, err := os.Executable()
	if err != nil {
		return "", "", fmt.Errorf("Cannot determine the path to the program, got %w", err)
	}
	var subCmd string
	if len(args) < 2 {
		subCmd = "status"
	} else {
		subCmd = args[1]
	}

	return mainCmd, subCmd, nil
}

func getVersion() error {
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
