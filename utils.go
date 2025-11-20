package otplock

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/mjwhitta/errors"
)

func execute(cmd string) (string, error) {
	var sh []string = []string{"bash", "-c"}

	if runtime.GOOS == "windows" {
		sh = []string{"cmd", "/C"}
	}

	return executeShell(sh, cmd)
}

func executeShell(sh []string, cmd string) (string, error) {
	var e error
	var o []byte

	sh = append(sh, cmd)

	//nolint:gosec // G204 - that's the whole point of this function
	if o, e = exec.Command(sh[0], sh[1:]...).Output(); e != nil {
		e = errors.Newf("command \"%s\" returned error:\n%w", cmd, e)
		return "", e
	}

	return strings.TrimSpace(string(o)), nil
}
