//go:build !darwin

package notification

func (n *Notifier) Show(_, _ string) error {
	return nil
}
