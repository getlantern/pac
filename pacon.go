// +build !darwin

package pacon

/*
#cgo linux pkg-config: gtk+-3.0
#cgo windows LDFLAGS: -l rasapi32 -l wininet

#include "pacon.h"
#include <stdlib.h>

const char* EMPTY_STRING = "";
*/
import "C"
import "unsafe"

/* Tells OS to configure proxy through `pacUrl` */
func PacOn(pacUrl string) {
	cPacUrl := C.CString(pacUrl)
	C.togglePac(C.PAC_ON, cPacUrl)
	C.free(unsafe.Pointer(cPacUrl))
}

/* Set proxy mode back to direct/none */
func PacOff() {
	C.togglePac(C.PAC_OFF, C.EMPTY_STRING)
}
