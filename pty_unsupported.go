//go:build !linux && !darwin

package execx

import (
	"errors"
	"os"
)

// ptyCheck keeps unsupported-platform failure deterministic and side-effect free.
func ptyCheck() error {
	return errors.New("execx: WithPTY is not supported on this platform")
}

// openPTY mirrors ptyCheck for callers that reach allocation directly.
func openPTY() (*os.File, *os.File, error) {
	return nil, nil, ptyCheck()
}
