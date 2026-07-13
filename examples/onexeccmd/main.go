package main

import (
	"github.com/goforj/execx"
	"os/exec"
	"syscall"
)

func main() {
	// OnExecCmd registers a callback to mutate the underlying exec.Cmd before start.

	// Example: exec cmd
	_, _ = execx.Command("printf", "hi").
		OnExecCmd(func(cmd *exec.Cmd) {
			cmd.SysProcAttr = &syscall.SysProcAttr{}
		}).
		Run()
}
