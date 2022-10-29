package cmd

import (
	"errors"
	"fmt"
	"os"

	"foxyshot/config"
)

var errUnknownSubCommand = errors.New("unknown subcommand")

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
		return newDefault().start(mainCmd, "run")
	case "stop":
		return newDefault().stop()
	case "status":
		return printStatus(newDefault())
	case "configure":
		return config.RunConfigure()
	case "version":
		return displayVersion()
	default:
		return errUnknownSubCommand
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
