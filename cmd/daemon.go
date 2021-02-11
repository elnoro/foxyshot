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

func startDaemon(mainCmd string) error {
	p, err := getPID()
	if err == nil {
		return fmt.Errorf("Daemon is already running, PID: %d", p)
	}
	cmd := exec.Command(mainCmd, "run")
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Cannot start daemon, got %v", err)
	}

	err = ioutil.WriteFile(getStateFile(), []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
	if err != nil {
		return fmt.Errorf("Cannot save the status of foxyshot daemon. PID: %d, error: %v",
			cmd.Process.Pid,
			err,
		)
	}

	return nil
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
