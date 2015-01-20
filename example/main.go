package main

import (
	"fmt"
	"github.com/fffw/proxysetup"
)

func main() {
	proxysetup.TurnOnAutoProxy("localhost:12345/pac")
	fmt.Println("proxy set, any key to continue...")
	var i int
	fmt.Scanf("%d\n", &i)
	proxysetup.TurnOffAutoProxy()
}
