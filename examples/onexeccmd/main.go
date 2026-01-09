//go:build ignore
// +build ignore

package main

import (
	"os/exec"

	"github.com/goforj/execx"
)

func main() {
	// OnExecCmd registers a callback to mutate the underlying exec.Cmd before start.

	// Example: exec cmd
	_, _ = execx.Command("printf", "hi").
		OnExecCmd(func(cmd *exec.Cmd) {
			cmd.Env = append(cmd.Env, "EXAMPLE=1")
		}).
		Run()
}
