//go:build unix

package execx

import (
	"os"
	"syscall"
)

// signalFromState exposes signal termination without treating it as an execution error.
func signalFromState(state *os.ProcessState) os.Signal {
	if state == nil {
		return nil
	}
	waitStatus := state.Sys().(syscall.WaitStatus)
	if waitStatus.Signaled() {
		return waitStatus.Signal()
	}
	return nil
}
