package main

import (
	"fmt"
	"runtime"

	"github.com/getlantern/pac"
)

func main() {
	if runtime.GOOS == "darwin" {
		pacon.SetHelperNameOnOSX("neat-helper")
		pacon.SetIconPathOnOSX("icon.png")
		pacon.SetPromptOnOSX("Input your password and save the world!")
	}
	err := pacon.PacOn("localhost:12345/pac")
	if err != nil {
		fmt.Printf("Error set proxy: %s\n", err)
		return
	}
	fmt.Println("proxy set, Enter continue...")
	var i int
	fmt.Scanf("%d\n", &i)
	pacon.PacOff()
}
