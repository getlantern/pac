package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/getlantern/pac"
)

func main() {
	pac.SetHelperPath(absPath("./pac-cmd"))
	if runtime.GOOS == "darwin" {
		pac.SetIconPathOnOSX(absPath("icon.png"))
		pac.SetPromptOnOSX("Input your password and save the world!")
	}
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

func absPath(name string) string {
	wd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return filepath.Join(wd, name)
}
