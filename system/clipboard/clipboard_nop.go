//go:build !darwin

package clipboard

func (m *Clipboard) Copy(_ string) error {
	return nil
}
