//go:build windows

package execx

import "os"

// terminateHelperProcess uses an exit code because Windows has no portable self-signal equivalent.
func terminateHelperProcess() {
	os.Exit(3)
}
