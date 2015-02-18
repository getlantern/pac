// +build !darwin

package pac

func ensureElevatedOnDarwin(path string) error {
	return nil
}
