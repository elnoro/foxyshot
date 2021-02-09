package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const defaultStateFile = "/tmp/foxyshot.state"

var stateFile = defaultStateFile

func getStateFile() string {
	return stateFile
}

func startDaemon(mainCmd string) {
	cmd := exec.Command(mainCmd, "run")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cannot start daemon, got %v", err)
	}

	// TODO fail gracefully if the state file exist
	err = ioutil.WriteFile(getStateFile(), []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
	if err != nil {
		log.Println(
			"Could save the status of foxyshot daemon. PID:",
			cmd.Process.Pid,
			"error:",
			err,
		)
	}
}

func stopDaemon() {
	pid, err := getPID()
	if err != nil {
		log.Printf("Cannot find the state of the app. Got %v\n", err)

		return
	}

	log.Println("Stopping process with pid", pid)

	cmd := exec.Command("kill", strconv.Itoa(pid))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Got error when stopping process: %v\n", err)
	}
	log.Println(out)
	os.Remove(getStateFile())
}

func printStatus() {
	pid, err := getPID()
	if err != nil {
		log.Printf("Cannot find the state of the app. Got %v\n", err)
		return
	}
	fmt.Println("Running. PID:", pid)
}

func getPID() (pid int, err error) {
	state, err := ioutil.ReadFile(getStateFile())
	if err == nil {
		fmt.Sscanf(string(state), "%d", &pid)
		if pid == 0 {
			return 0, fmt.Errorf("unexpected pid value %d", pid)
		}
	}
	return
}
