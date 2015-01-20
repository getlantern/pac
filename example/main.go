package main

import (
	"github.com/getlantern/proxysetup"
)

func main() {
	proxysetup.TurnOnAutoProxy("localhost:12345/pac")
	proxysetup.TurnOffAutoProxy()
}
