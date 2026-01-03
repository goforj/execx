//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Setpgid is a no-op on Windows; on Unix it places the child in a new process group.

	// Example: setpgid
	out, _ := execx.Command("printf", "ok").Setpgid(true).Output()
	fmt.Print(out)
	// ok
}
