//go:build darwin

package clipboard

import (
	"fmt"
	"io"
	"os/exec"
)

// Copy uses pbcopy to copy values into the system clipboard (macos only)
func (m *Clipboard) Copy(val string) error {
	cmd := exec.Command("pbcopy")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("clipboard error, %w", err)
	}
	_, err = io.WriteString(stdin, val)
	if err != nil {
		return fmt.Errorf("clipboard error, %w", err)
	}
	err = stdin.Close()
	if err != nil {
		return fmt.Errorf("clipboard error, %w", err)
	}
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("clipboard error, %w", err)
	}

	return nil
}
