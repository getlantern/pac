// +build !darwin

package pac

func SetIconPathOnMacOS(i string) {
}

func SetPromptOnMacOS(p string) {
}

func prestine(path string) bool {
	return true
}

func elevateOnDarwin(path string) error {
	return nil
}
