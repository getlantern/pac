// +build darwin

package pacon

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework SystemConfiguration -framework Security

#include "pacon.h"
#include <stdlib.h>

const char* EMPTY_STRING = "";
int runAuthorized(const char *path);
int togglePacWithHelper(int onOff, const char* autoProxyConfigFileUrl, const char* path);
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

// On Mac OSX, we need a previledged action.
// SetHelperNameOnOSX sets the file name to generated.
func SetHelperNameOnOSX(name string) {
	helperToolName = name
}

func helperAbsPath() string {
	wd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return filepath.Join(wd, helperToolName)
}

func doTogglePac(onOff C.int, pacUrl *C.char) (err error) {
	if err = ensureHelperTool(); err != nil {
		err = fmt.Errorf("Unable to extract helper tool: %s", err)
		return
	}
	ret := C.togglePacWithHelper(onOff, pacUrl, C.CString(helperAbsPath()))
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
	return doTogglePac(C.PAC_OFF, C.EMPTY_STRING)
}

func ensureHelperTool() (err error) {
	absPath := helperAbsPath()
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
	ret := C.runAuthorized(C.CString(path))
	if ret != 0 {
		err = fmt.Errorf("Unable to runAuthorized on helper tool")
	}
	return
}
