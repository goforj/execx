//go:build !unix && !windows

package execx

import "os"

// terminateHelperProcess uses an exit code where no portable self-signal mechanism exists.
func terminateHelperProcess() {
	os.Exit(3)
}
