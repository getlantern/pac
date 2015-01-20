package proxysetup

/*
#cgo linux pkg-config: gtk+-3.0
#cgo linux CFLAGS: -DLINUX
#cgo windows CFLAGS: -DWIN32 -Wl, -l rasapi32 -Wl, -l wininet
#cgo windows LDFLAGS: -l rasapi32 -l wininet
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework SystemConfiguration -framework Security

#include "proxysetup.h"
*/
import "C"

func TurnOnAutoProxy(pacUrl string) {
	C.toggleAutoProxyConfigFile(C.CString("on"), C.CString(pacUrl))
}

func TurnOffAutoProxy() {
	C.toggleAutoProxyConfigFile(C.CString("off"), C.CString(""))
}
