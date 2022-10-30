//go:build !darwin

package clipboard

func (m *Clipboard) Copy(val string) error {
	return nil
}
