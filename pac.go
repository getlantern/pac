package pac

import "C"
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var helperToolName string = "pac-cmd"

// This library will extract a helper tool under application's same directory
// to actually change proxy setup.
// SetHelperName specifies the file name to be generated.
func SetHelperName(name string) {
	helperToolName = name
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

func run(cmd *exec.Cmd) error {
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Unable to execute pac tool: %s\n%s", err, out.String())
	}
	return nil
}

func ensureHelperTool() (err error) {
	absPath := absPath(helperToolName)
	if _, err = os.Stat(absPath); err != nil {
		err = extractHelper(absPath)
	} else if !prestine(absPath) {
		os.Remove(absPath)
		if err != nil {
			err = fmt.Errorf("Error remove existing %s: %s", absPath, err)
		}
		err = extractHelper(absPath)
	}
	return
}

func extractHelper(path string) error {
	err := ioutil.WriteFile(path, pacBytes, 0755)
	if err != nil {
		return fmt.Errorf("Error write helper file %s: %s", path, err)
	}
	return elevateOnDarwin(path)
}

func absPath(name string) string {
	wd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return filepath.Join(wd, name)
}


