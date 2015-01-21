package pacon

/*
#cgo linux pkg-config: gtk+-3.0
#cgo linux CFLAGS: -DLINUX
#cgo windows CFLAGS: -DWIN32
#cgo windows LDFLAGS: -l rasapi32 -l wininet
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework SystemConfiguration -framework Security

#include "pacon.h"
*/
import "C"

func PacOn(pacUrl string) {
	C.togglePac(C.PAC_ON, C.CString(pacUrl))
}

func PacOff() {
	C.togglePac(C.PAC_OFF, C.CString(""))
}
