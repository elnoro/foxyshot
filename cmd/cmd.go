package cmd

import (
	"errors"
	"fmt"
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

	case "help":
		fallthrough
	case "-h":
		return printHelp()

	case "version":
		fallthrough
	case "-v":
		fallthrough
	case "--version":
		return printVersion()
	default:
		return errUnknownSubCommand
	}
}

func printHelp() error {
	fmt.Print(`
🦊FoxyShot is a lightweight tool to upload MacOS screenshots to an S3-compatible provider🦊
Usage: foxyshot [command]
Available commands:
	  run        Run foxyshot in foreground
	  configure  Configure foxyshot

	  start      Start foxyshot daemon
	  stop       Stop foxyshot daemon
	  status     Print status of foxyshot daemon

	  help       Print this help message
	  version    Print version
Available flags:
	  -logfile		  Path to the log file (default: STDOUT)
`)

	return nil
}

func parseArgs(args []string) string {
	if len(args) < 2 {
		return "help"
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
