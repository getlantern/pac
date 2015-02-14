package pac

import "C"
import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

var helperToolName string = "helper"

func absPath(name string) string {
	wd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return filepath.Join(wd, name)
}

/* On tells OS to configure proxy through `pacUrl` */
func On(pacUrl string) (err error) {
	if err = ensureHelperTool(); err != nil {
		err = fmt.Errorf("Unable to extract helper tool: %s", err)
		return
	}
	cmd := exec.Command(absPath(helperToolName), "on", pacUrl)
	return run(cmd)
}

/* Off sets proxy mode back to direct/none */
func Off() (err error) {
	if err = ensureHelperTool(); err != nil {
		err = fmt.Errorf("Unable to extract helper tool: %s", err)
		return
	}
	cmd := exec.Command(absPath(helperToolName), "off")
	return run(cmd)
}

func run(cmd exec.Cmd) error {
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Unable to execute pac tool: %s\n%s", err, out.String())
	}
}

func ensureHelperTool() (err error) {
	absPath := absPath(helperToolName)
	var s syscall.Stat_t
	err = syscall.Stat(absPath, &s)
	if err != nil {
		err = extractHelper(absPath)
	} else if s.Mode&syscall.S_ISUID == 0 || s.Uid != 0 || s.Gid != 0 {
		os.Remove(absPath)
		if err != nil {
			err = fmt.Errorf("Error remove existing %s: %s", absPath, err)
		}
		err = extractHelper(absPath)
	}
	return
}

func extractHelper(path string) (err error) {
	err = ioutil.WriteFile(path, helper, syscall.S_IRWXU)
	if err != nil {
		err = fmt.Errorf("Error write helper file %s: %s", path, err)
	}
	if runtime.GOOS == "darwin" {
		err = elevate(path)
	}
	return
}
