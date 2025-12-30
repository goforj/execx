package execx

import (
	"os"
	"time"
)

// Result captures the outcome of a command execution.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
	Duration time.Duration

	signal os.Signal
}

// OK reports whether the command exited cleanly without errors.
// @group Results
//
// Example: ok
//
//	res, err := execx.Command("go", "env", "GOOS").Run()
//	fmt.Println(err == nil && res.OK())
//	// #bool true
func (r Result) OK() bool {
	return r.Err == nil && r.ExitCode == 0
}

// IsExitCode reports whether the exit code matches.
// @group Results
//
// Example: exit code
//
//	res, err := execx.Command("go", "env", "GOOS").Run()
//	fmt.Println(err == nil && res.IsExitCode(0))
//	// #bool true
func (r Result) IsExitCode(code int) bool {
	return r.ExitCode == code
}

// IsSignal reports whether the command terminated due to a signal.
// @group Results
//
// Example: signal
//
//	res, err := execx.Command("go", "env", "GOOS").Run()
//	fmt.Println(err == nil && res.IsSignal(os.Interrupt))
//	// #bool false
func (r Result) IsSignal(sig os.Signal) bool {
	return r.signal == sig
}
