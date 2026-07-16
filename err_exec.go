package execx

import "os"

// ErrExec reports a failure to start or an explicit execution failure.
type ErrExec struct {
	// Err contains the underlying start or I/O failure.
	Err error
	// ExitCode contains the child status when a process reached execution.
	ExitCode int
	// Signal contains the terminating signal on platforms that expose one.
	Signal os.Signal
	// Stderr contains captured diagnostic output available when the failure occurred.
	Stderr string
}

// Error returns the wrapped error message when available.
// @group Errors
//
// Example: error string
//
//	err := execx.ErrExec{Err: fmt.Errorf("boom")}
//	fmt.Println(err.Error())
//	// #string boom
func (e ErrExec) Error() string {
	if e.Err == nil {
		return "execx: execution failed"
	}
	return e.Err.Error()
}

// Unwrap exposes the underlying error.
// @group Errors
//
// Example: unwrap
//
//	err := execx.ErrExec{Err: fmt.Errorf("boom")}
//	fmt.Println(err.Unwrap() != nil)
//	// #bool true
func (e ErrExec) Unwrap() error {
	return e.Err
}
