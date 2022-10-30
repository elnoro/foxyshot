//go:build !darwin

package notification

func (n *Notifier) Show(title, notification string) error {
	return nil
}
