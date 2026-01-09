//go:build !linux && !darwin

package execx

import (
	"errors"
	"os"
)

func ptyCheck() error {
	return errors.New("execx: WithPTY is not supported on this platform")
}

func openPTY() (*os.File, *os.File, error) {
	return nil, nil, ptyCheck()
}
