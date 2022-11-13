package cmd

import (
	"errors"
	"log"
	"os"

	"foxyshot/config"
	"foxyshot/system/logger"
)

var errUnknownSubCommand = errors.New("unknown subcommand")

// RunCmd parses the subcommand and chooses the behaviour
func RunCmd(args []string) error {
	subCmd := parseArgs(args)
	switch subCmd {
	case "run":
		return run()
	case "start":
		return newDefaultDaemon().start(getExecutable(), "run")
	case "stop":
		return newDefaultDaemon().stop()
	case "status":
		return printStatus(newDefaultDaemon())
	case "configure":
		return config.RunConfigure()
	case "version":
		return printVersion()
	default:
		return errUnknownSubCommand
	}
}

func parseArgs(args []string) string {
	if len(args) < 2 {
		return "status"
	}
	subCmd := args[1]
	logger.FromArgs(args[2:])

	return subCmd
}

func getExecutable() string {
	path, err := os.Executable()
	if err != nil {
		log.Fatalf("Cannot determine the path to the program, got %v", err)
	}
	return path
}
