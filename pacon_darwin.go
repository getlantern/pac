// +build darwin

package pacon

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework SystemConfiguration -framework Security

#include "pacon.h"
#include <stdlib.h>

const char* EMPTY_STRING = "";
void runAuthorized(const char *path);
void togglePacWithHelper(int onOff, const char* autoProxyConfigFileUrl, const char* path);
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

var helperToolName string = "helper"

/* Tells OS to configure proxy through `pacUrl` */
func PacOn(pacUrl string) {
	ensureHelperTool()
	C.togglePacWithHelper(C.PAC_ON, C.CString(pacUrl), C.CString(helperToolName))
}

/* Set proxy mode back to direct/none */
func PacOff() {
	ensureHelperTool()
	C.togglePacWithHelper(C.PAC_OFF, C.EMPTY_STRING, C.CString(helperToolName))
}

func ensureHelperTool() {
	var s syscall.Stat_t
	err := syscall.Stat(helperToolName, &s)
	if err != nil {
		fmt.Printf("%v\n", err)
		extractHelper()
	} else if s.Mode&syscall.S_ISUID == 0 || s.Uid != 0 || s.Gid != 0 {
		fmt.Printf("%v %v\n", s.Mode, s.Uid)
		os.Remove(helperToolName)
		extractHelper()
	}
}

func extractHelper() {
	err := ioutil.WriteFile(helperToolName, helper, syscall.S_IRWXU)
	if err != nil {
		return
	}
	C.runAuthorized(C.CString(helperToolName))
}
