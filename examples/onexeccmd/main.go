package main

import (
	"github.com/goforj/execx"
	"os/exec"
	"syscall"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// OnExecCmd registers a callback to mutate the underlying exec.Cmd before execx
	// attaches its stdin, output capture, and pipeline wiring. Use it for fields such
	// as SysProcAttr; execx owns Stdin, Stdout, and Stderr during execution.

	// Example: exec cmd
	_, _ = execx.Command("printf", "hi").
		OnExecCmd(func(cmd *exec.Cmd) {
			cmd.SysProcAttr = &syscall.SysProcAttr{}
		}).
		Run()
}
