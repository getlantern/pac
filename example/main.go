package main

import (
	"fmt"

	"github.com/getlantern/pacon"
)

func main() {
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
