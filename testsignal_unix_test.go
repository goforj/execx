//go:build unix

package execx

import (
	"os"
	"syscall"
)

// terminateHelperProcess delivers a real signal so Unix result mapping is exercised end to end.
func terminateHelperProcess() {
	_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
}
