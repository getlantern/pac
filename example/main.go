package main

import (
	"fmt"

	"github.com/getlantern/pac"
)

func main() {
	pac.SetHelperPath("./pac-cmd")
	pac.SetIconPathOnMacOS("./icon.png")
	pac.SetPromptOnMacOS("Input your password and save the world!")
	err := pac.On("localhost:12345/pac")
	if err != nil {
		fmt.Printf("Error set proxy: %s\n", err)
		return
	}
	fmt.Println("proxy set, Enter continue...")
	var i int
	fmt.Scanf("%d\n", &i)
	pac.Off()
}
