package pac

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"

	"github.com/getlantern/byteexec"
)

var mu sync.Mutex
var be *byteexec.Exec

// EnsureHelperToolPresent checks if helper tool exists and extracts it if not.
// On Mac OS, it also checks and set the file's owner to root:wheel and the setuid bit,
// it will request user to input password through a dialog to gain the rights to do so.
// fullPath: the file to be checked and generated if not exists.
// prompt: the message to be shown on the dialog.
// iconPath: the full path of the icon to be shown on the dialog.
func EnsureHelperToolPresent(fullPath string, prompt string, iconFullPath string) (err error) {
	mu.Lock()
	defer mu.Unlock()
	be, err = byteexec.New(pacBytes, fullPath)
	if err != nil {
		return fmt.Errorf("Unable to extract helper tool: %s", err)
	}
	return ensureElevatedOnDarwin(fullPath, prompt, iconFullPath)
}

/* On tells OS to configure proxy through `pacUrl` */
func On(pacUrl string) (err error) {
	mu.Lock()
	defer mu.Unlock()
	if be == nil {
		return fmt.Errorf("call EnsureHelperToolPresent() first")
	}
	cmd := be.Command("on", pacUrl)
	return run(cmd)
}

/* Off sets proxy mode back to direct/none */
func Off() (err error) {
	mu.Lock()
	defer mu.Unlock()
	if be == nil {
		return fmt.Errorf("call EnsureHelperToolPresent() first")
	}
	cmd := be.Command("off")
	return run(cmd)
}

func run(cmd *exec.Cmd) error {
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Unable to execute pac tool: %s\n%s", err, errOut.String())
	}
	return nil
}
