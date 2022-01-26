package cmd

import (
	"fmt"
	"foxyshot/config"
	"os"
)

// RunCmd parses the subcommand and chooses the behaviour
func RunCmd(args []string) error {
	mainCmd, subCmd, err := parseArgs(args)
	if err != nil {
		return err
	}
	switch subCmd {
	case "run":
		return startApp()
	case "start":
		return startDaemon(mainCmd)
	case "stop":
		return stopDaemon()
	case "status":
		return printStatus()
	case "configure":
		return config.RunConfigure()
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
