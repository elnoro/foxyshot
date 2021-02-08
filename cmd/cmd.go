package cmd

import (
	"log"
)

// RunCmd parses the subcommand and chooses the behaviour
func RunCmd(args []string) {
	mainCmd, subCmd := parseArgs(args)
	switch subCmd {
	case "run":
		startApp()
	case "start":
		startDaemon(mainCmd)
	case "stop":
		stopDaemon()
	case "status":
		printStatus()
	default:
		log.Println("Unknown subcommand:", subCmd)
	}
}

func parseArgs(args []string) (mainCmd string, subCmd string) {
	mainCmd = args[0]
	if len(args) < 2 {
		subCmd = "status"
	} else {
		subCmd = args[1]
	}

	return
}
