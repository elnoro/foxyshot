package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

const defaultStateFile = "/tmp/foxyshot.state"

type daemon struct {
	stateFile string
}

func newDefaultDaemon() *daemon {
	return newDaemon(defaultStateFile)
}

func newDaemon(stateFile string) *daemon {
	return &daemon{stateFile: stateFile}
}

func (d *daemon) start(args ...string) error {
	p, err := d.getPID()
	if err == nil {
		return fmt.Errorf("Daemon is already running, PID: %d, reason %w", p, err)
	}
	cmd := exec.Command(args[0], args[1:]...)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Cannot start daemon, reason %w", err)
	}

	err = os.WriteFile(d.stateFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0600)
	if err != nil {
		return fmt.Errorf("Cannot save the status of foxyshot daemon. PID: %d, error: %w",
			cmd.Process.Pid,
			err,
		)
	}

	return nil
}

func (d *daemon) stop() error {
	pid, err := d.getPID()
	if err != nil {
		return fmt.Errorf("Cannot find the state of the app. Got %w", err)
	}
	if pid == 0 {
		return fmt.Errorf("Invalid pid, cannot be zero. Check the state file")
	}

	log.Println("Stopping process with pid", pid)

	err = syscall.Kill(pid, syscall.SIGINT)
	if errors.Is(err, syscall.ESRCH) {
		log.Println("Process is not running. Removing state")
	} else if err != nil {
		return fmt.Errorf("Got error when stopping process: %w", err)
	}
	err = os.Remove(d.stateFile)
	if err != nil {
		return fmt.Errorf("Got error when removing state file: %w", err)
	}

	return nil
}

func (d *daemon) getPID() (pid int, err error) {
	state, err := os.ReadFile(d.stateFile)
	if err != nil {
		return 0, fmt.Errorf("getting pid error, reason %w", err)
	}
	_, err = fmt.Sscanf(string(state), "%d", &pid)
	if pid == 0 {
		return 0, fmt.Errorf("unexpected pid value %d, reason %w", pid, err)
	}
	return
}

func printStatus(d *daemon) error {
	pid, err := d.getPID()
	if err != nil {
		return fmt.Errorf("Printing status error, %w", err)
	}
	fmt.Println("Running. PID:", pid)

	return nil
}
