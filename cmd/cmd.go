package cmd

import (
	"foxyshot/config"
	"log"
	"os"
)

// RunCmd parses the subcommand and chooses the behaviour
func RunCmd(args []string) {
	mainCmd, subCmd := parseArgs(args)
	switch subCmd {
	case "run":
		startApp()
	case "start":
		wrapError(startDaemon(mainCmd))
	case "stop":
		stopDaemon()
	case "status":
		printStatus()
	case "configure":
		wrapError(config.RunConfigure())
	default:
		log.Println("Unknown subcommand:", subCmd)
	}
}

func parseArgs(args []string) (mainCmd string, subCmd string) {
	mainCmd, err := os.Executable()
	if err != nil {
		log.Fatal("Cannot determine the path to the program, got", err)
	}
	if len(args) < 2 {
		subCmd = "status"
	} else {
		subCmd = args[1]
	}

	return
}

func wrapError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
