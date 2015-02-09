// +build darwin

package pacon

/*
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework SystemConfiguration -framework Security

#include "pacon.h"
#include <stdlib.h>

const int PAC_ON = 1;
const int PAC_OFF = 0;

const char* NULL_STRING = NULL;
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

var helperToolName string = "helper"
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

func absPath(name string) string {
	wd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return filepath.Join(wd, name)
}

func doTogglePac(onOff C.int, pacUrl *C.char) (err error) {
	if err = ensureHelperTool(); err != nil {
		err = fmt.Errorf("Unable to extract helper tool: %s", err)
		return
	}
	ret := C.togglePacWithHelper(onOff, pacUrl, C.CString(absPath(helperToolName)))
	if ret != 0 {
		err = fmt.Errorf("Failed to run helper tool to set pac")
	}
	return
}

/* PacOn tells OS to configure proxy through `pacUrl` */
func PacOn(pacUrl string) (err error) {
	cPacUrl := C.CString(pacUrl)
	defer C.free(unsafe.Pointer(cPacUrl))
	return doTogglePac(C.PAC_ON, cPacUrl)
}

/* PacOff sets proxy mode back to direct/none */
func PacOff() (err error) {
	return doTogglePac(C.PAC_OFF, C.NULL_STRING)
}

func ensureHelperTool() (err error) {
	absPath := absPath(helperToolName)
	var s syscall.Stat_t
	err = syscall.Stat(absPath, &s)
	if err != nil {
		err = extractHelper(absPath)
	} else if s.Mode&syscall.S_ISUID == 0 || s.Uid != 0 || s.Gid != 0 {
		os.Remove(absPath)
		if err != nil {
			err = fmt.Errorf("Error remove existing %s: %s", absPath, err)
		}
		err = extractHelper(absPath)
	}
	return
}

func extractHelper(path string) (err error) {
	err = ioutil.WriteFile(path, helper, syscall.S_IRWXU)
	if err != nil {
		err = fmt.Errorf("Error write helper file %s: %s", path, err)
	}
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
