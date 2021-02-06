package main

import (
	"context"
	"fmt"
	"foxyshot/app"
	"foxyshot/config"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	mainCmd, subCmd := parseArgs(os.Args)
	switch subCmd {
	case "run":
		startApp()
		break
	case "start":
		startDaemon(mainCmd)
		break
	case "stop":
		stopDaemon()
		break
	default:
		log.Println("Unknown subcommand:", subCmd)
	}
}

func parseArgs(args []string) (mainCmd string, subCmd string) {
	mainCmd = os.Args[0]
	if len(args) < 2 {
		subCmd = "run"
	} else {
		subCmd = os.Args[1]
	}

	return
}

func startApp() {
	appConfig := config.Load()
	app, err := app.New(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go app.Watch(ctx, appConfig.WatchFor)

	app.WaitForExit(cancel)
}

const appStateFile = "/tmp/foxyshot.state"

func startDaemon(mainCmd string) {
	cmd := exec.Command(mainCmd, "run")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cannot start daemon, got %v", err)
	}

	// TODO fail gracefully if the state file exist
	ioutil.WriteFile(appStateFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
}

func stopDaemon() {
	state, err := ioutil.ReadFile(appStateFile)
	if err != nil {
		log.Printf("Cannot find the state of the app. Got %v\n", err)

		return
	}
	os.Remove(appStateFile)

	var pid int
	fmt.Sscanf(string(state), "%d", &pid)
	log.Println("Stopping process with pid", pid)

	cmd := exec.Command("kill", strconv.Itoa(pid))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Got error when stopping process: %v\n", err)
	}
	log.Println(out)
}
