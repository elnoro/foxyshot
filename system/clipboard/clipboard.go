package clipboard

import (
	"fmt"
	"io"
	"os/exec"
)

func NewClipboard() Clipboard {
	return &macosClipboard{}
}

type Clipboard interface {
	Copy(val string) error
}

type macosClipboard struct {
}

// Copy uses pbcopy to copy values into the system clipboard (macos only)
func (m macosClipboard) Copy(val string) error {
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
