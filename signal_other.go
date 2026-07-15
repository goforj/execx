//go:build !unix && !windows

package execx

import "os"

// signalFromState returns nil because these platforms expose no portable process-signal mapping.
func signalFromState(_ *os.ProcessState) os.Signal {
	return nil
}
