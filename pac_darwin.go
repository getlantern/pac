package pac

import "C"

/*
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework SystemConfiguration -framework Security
#include "darwin.h"

const char* NULL_STRING = "";
*/

var iconPath string
var prompt string

// On Mac OSX, we'll extract a helper tool with root priviledge
// under application's same directory to actually change proxy setup,
// SetHelperNameOnOSX specifies the file name to be generated.
func SetHelperNameOnOSX(name string) {
	helperToolName = name
}

// Mac OSX will show a dialog requesting user to input password,
// SetIconPathOnOSX specifies the icon to be shown on the dialog.
func SetIconPathOnOSX(i string) {
	iconPath = i
}

// Mac OSX will show a dialog requesting user to input password,
// SetPromptOnOSX specifies the text to be shown on this dialog.
func SetPromptOnOSX(p string) {
	prompt = p
}

func prestine(path string) bool {
	var s syscall.Stat_t
	// we just checked its existence, not bother checking specific error again
	if err = syscall.Stat(absPath, &s); err != nil {
		return false;
	}
	if s.Mode&syscall.S_ISUID == 0 || s.Uid != 0 || s.Gid != 0 {
		return false;
	}
	return true;
}

func elevateOnDarwin(path string) (err error) {
	cPrompt := C.NULL_STRING
	if prompt != "" {
		cPrompt = C.CString(prompt)
	}
	cIconPath := C.NULL_STRING
	if iconPath != "" {
		if !filepath.IsAbs(iconPath) {
			iconPath = absPath(iconPath)
		}
		cIconPath = C.CString(iconPath)
	}
	ret := C.runAuthorized(C.CString(path), cPrompt, cIconPath)
	if ret != 0 {
		err = fmt.Errorf("Unable to runAuthorized on helper tool")
	}
	return
}
