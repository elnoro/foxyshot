package cmd

import (
	"fmt"
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
		return fmt.Errorf("Daemon is already running, PID: %d, reason %w", p, err)
	}
	cmd := exec.Command(mainCmd, "run")
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Cannot start daemon, reason %w", err)
	}

	err = os.WriteFile(getStateFile(), []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
	if err != nil {
		return fmt.Errorf("Cannot save the status of foxyshot daemon. PID: %d, error: %w",
			cmd.Process.Pid,
			err,
		)
	}

	return nil
}

func stopDaemon() error {
	pid, err := getPID()
	if err != nil {
		return fmt.Errorf("Cannot find the state of the app. Got %w\n", err)
	}

	log.Println("Stopping process with pid", pid)

	cmd := exec.Command("kill", strconv.Itoa(pid))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Got error when stopping process: %w\n", err)
	}
	log.Println(out)
	err = os.Remove(getStateFile())
	if err != nil {
		return fmt.Errorf("Got error when stopping process: %v\n", err)
	}

	return nil
}

func printStatus() error {
	pid, err := getPID()
	if err != nil {
		return fmt.Errorf("Printing status error, %w", err)
	}
	fmt.Println("Running. PID:", pid)

	return nil
}

func getPID() (pid int, err error) {
	state, err := os.ReadFile(getStateFile())
	if err != nil {
		return 0, fmt.Errorf("getting pid error, reason %w", err)
	}
	_, err = fmt.Sscanf(string(state), "%d", &pid)
	if pid == 0 {
		return 0, fmt.Errorf("unexpected pid value %d, reason %w", pid, err)
	}
	return
}
