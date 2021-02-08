package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

// TODO needs to be configurable, at least for testing
const appStateFile = "/tmp/foxyshot.state"

func startDaemon(mainCmd string) {
	cmd := exec.Command(mainCmd, "run")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cannot start daemon, got %v", err)
	}

	// TODO fail gracefully if the state file exist
	err = ioutil.WriteFile(appStateFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
	if err != nil {
		log.Println("Could save the status of foxyshot daemon. PID: ", cmd.Process.Pid)
	}
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

func printStatus() {
	state, err := ioutil.ReadFile(appStateFile)
	if err != nil {
		log.Printf("Cannot find the state of the app. Got %v\n", err)

		return
	}

	var pid int
	fmt.Sscanf(string(state), "%d", &pid)
	fmt.Println("Running. PID:", pid)
}
