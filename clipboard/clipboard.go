package clipboard

import (
	"io"
	"os/exec"
)

// CopyToClipboard uses pbcopy to copy values into the system clipboard (macos only)
func CopyToClipboard(val string) error {
	cmd := exec.Command("pbcopy")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	_, err = io.WriteString(stdin, val)
	if err != nil {
		return err
	}
	err = stdin.Close()
	if err != nil {
		return err
	}

	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
