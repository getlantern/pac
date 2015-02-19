// +build !darwin

package pac

func ensureElevatedOnDarwin(be *byteexec.Exec, helperFullPath string, prompt string, iconFullPath string) (err error) {
	return nil
}
