package notification

import (
	"fmt"
	"os/exec"
)

func NewNotifier() Notifier {
	return &macosNotifier{}
}

type macosNotifier struct {
}

func (m *macosNotifier) Show(title, notification string) error {
	processPath, err := exec.LookPath("osascript")
	if err != nil {
		return fmt.Errorf("notification error, %w", err)
	}

	osaCommand := fmt.Sprintf("display notification \"%s\" with title \"%s\"", notification, title)
	cmd := exec.Command(processPath, "-e", osaCommand)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("notification error, %w", err)
	}

	return nil
}

type Notifier interface {
	Show(title, notification string) error
}
