//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// CreationFlags sets Windows process creation flags (for example, create a new process group).

	// Example: creation flags
	out, _ := execx.Command("printf", "ok").CreationFlags(execx.CreateNewProcessGroup).Output()
	fmt.Print(out)
	// ok
}
