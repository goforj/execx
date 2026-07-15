//go:build windows

package execx

import "os"

// signalFromState returns nil because Windows ProcessState has no portable os.Signal mapping.
func signalFromState(_ *os.ProcessState) os.Signal {
	return nil
}
