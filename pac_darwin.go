package pac

import (
	"github.com/getlantern/byteexec"
	"github.com/getlantern/elevate"
)

import (
	"fmt"
	"syscall"
)

func ensureElevatedOnDarwin(be *byteexec.Exec, helperFullPath string, prompt string, iconFullPath string) (err error) {
	var s syscall.Stat_t
	// we just checked its existence, not bother checking specific error again
	if err = syscall.Stat(helperFullPath, &s); err != nil {
		return fmt.Errorf("Error stating helper tool %s: %s", err)
	}
	if s.Mode&syscall.S_ISUID > 0 && s.Uid == 0 && s.Gid == 0 {
		return
	}
	cmd := elevate.PromptWithIcon(prompt, iconFullPath, be.Filename, "setuid")
	return run(cmd)
}
