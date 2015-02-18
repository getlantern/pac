package pac

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var helperPath string = "./pac"

// This library will extract a helper tool to actually change pac.
// SetHelperPath specifies the file path to be generated.
// It will be 'pac' under current work directory by default.
func SetHelperPath(path string) {
	helperPath = path
}

/* On tells OS to configure proxy through `pacUrl` */
func On(pacUrl string) (err error) {
	if err = ensureHelperTool(); err != nil {
		err = fmt.Errorf("Unable to extract helper tool: %s", err)
		return
	}
	cmd := exec.Command(helperPath, "on", pacUrl)
	return run(cmd)
}

/* Off sets proxy mode back to direct/none */
func Off() (err error) {
	if err = ensureHelperTool(); err != nil {
		err = fmt.Errorf("Unable to extract helper tool: %s", err)
		return
	}
	cmd := exec.Command(helperPath, "off")
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
	if _, err = os.Stat(helperPath); err != nil {
		err = extractHelper(helperPath)
	} else if !prestine(helperPath) {
		// remove first so we can write even if we don't have written permission to override directly
		os.Remove(helperPath)
		if err != nil {
			err = fmt.Errorf("Error remove existing %s: %s", helperPath, err)
		}
		err = extractHelper(helperPath)
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
