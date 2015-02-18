package pac

/*
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework Security

extern int runAuthorized(char *path, char *prompt, char *iconPath);
const char* null() { return 0; }
*/
import "C"
import (
	"fmt"
	"syscall"
)

var iconPath string
var prompt string

// SetIconPathOnMacOS specifies the icon to be shown on the dialog.
func SetIconPathOnMacOS(i string) {
	iconPath = i
}

// SetPromptOnMacOS specifies the text to be shown on the dialog.
func SetPromptOnMacOS(p string) {
	prompt = p
}

func prestine(path string) bool {
	var s syscall.Stat_t
	// we just checked its existence, not bother checking specific error again
	if err := syscall.Stat(path, &s); err != nil {
		return false
	}
	if s.Mode&syscall.S_ISUID == 0 || s.Uid != 0 || s.Gid != 0 {
		return false
	}
	return true
}

func elevateOnDarwin(path string) (err error) {
	cPrompt := C.null()
	if prompt != "" {
		cPrompt = C.CString(prompt)
	}
	cIconPath := C.null()
	if iconPath != "" {
		cIconPath = C.CString(iconPath)
	}
	ret := C.runAuthorized(C.CString(path), cPrompt, cIconPath)
	if ret != 0 {
		err = fmt.Errorf("Unable to runAuthorized on helper tool")
	}
	return
}
