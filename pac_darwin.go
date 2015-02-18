package pac

/*
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework Security

extern int runAuthorized(char *helperFullPath, char *prompt, char *iconFullPath);
const char* null() { return 0; }
*/
import "C"
import (
	"fmt"
	"syscall"
)

func ensureElevatedOnDarwin(helperFullPath string, prompt string, iconFullPath string) (err error) {
	var s syscall.Stat_t
	// we just checked its existence, not bother checking specific error again
	if err = syscall.Stat(helperFullPath, &s); err != nil {
		return fmt.Errorf("Error stating helper tool %s: %s", err)
	}
	if s.Mode&syscall.S_ISUID > 0 && s.Uid == 0 && s.Gid == 0 {
		return
	}
	cPrompt := C.null()
	if prompt != "" {
		cPrompt = C.CString(prompt)
	}
	cIconFullPath := C.null()
	if iconFullPath != "" {
		cIconFullPath = C.CString(iconFullPath)
	}
	ret := C.runAuthorized(C.CString(helperFullPath), cPrompt, cIconFullPath)
	if ret != 0 {
		return fmt.Errorf("Unable to runAuthorized on helper tool")
	}
	return
}
