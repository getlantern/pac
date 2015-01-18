package main

import (
	"github.com/getlantern/proxysetup"
)

func main() {
	proxysetup.TurnOnAutoProxy("a.com")
	proxysetup.TurnOffAutoProxy()
}
