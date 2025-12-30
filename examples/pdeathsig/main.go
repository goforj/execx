//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"syscall"
)

func main() {
	// Pdeathsig is a no-op on Windows; on Linux it signals the child when the parent exits.

	// Example: pdeathsig
	out, _ := execx.Command("printf", "ok").Pdeathsig(syscall.SIGTERM).Output()
	fmt.Print(out)
	// ok
}
