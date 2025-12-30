package execx

import "os"

// ErrExec reports a failure to start or an explicit execution failure.
type ErrExec struct {
	Err      error
	ExitCode int
	Signal   os.Signal
	Stderr   string
}

func (e ErrExec) Error() string {
	if e.Err == nil {
		return "execx: execution failed"
	}
	return e.Err.Error()
}

func (e ErrExec) Unwrap() error {
	return e.Err
}
