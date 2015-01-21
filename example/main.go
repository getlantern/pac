package main

import (
	"fmt"

	"github.com/getlantern/pacon"
)

func main() {
	pacon.PacOn("localhost:12345/pac")
	fmt.Println("proxy set, any key to continue...")
	var i int
	fmt.Scanf("%d\n", &i)
	pacon.PacOff()
}
