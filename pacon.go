// +build !darwin

package pacon

/*
#cgo linux pkg-config: gtk+-3.0
#cgo windows LDFLAGS: -l rasapi32 -l wininet

#include "pacon.h"
#include <stdlib.h>

const int PAC_ON = 1;
const int PAC_OFF = 0;
const char* NULL_STRING = NULL;
*/
import "C"
import "unsafe"

func SetHelperNameOnOSX(name string) {
}

func SetIconPathOnOSX(i string) {
}

func SetPromptOnOSX(p string) {
}

/* Tells OS to configure proxy through `pacUrl` */
func PacOn(pacUrl string) (err error) {
	cPacUrl := C.CString(pacUrl)
	C.togglePac(C.PAC_ON, cPacUrl)
	C.free(unsafe.Pointer(cPacUrl))
	return
}

/* Set proxy mode back to direct/none */
func PacOff() (err error) {
	C.togglePac(C.PAC_OFF, C.NULL_STRING)
	return
}
