package main

import (
	"fmt"
	"github.com/goforj/execx"
	"syscall"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Pdeathsig is a no-op on non-Linux platforms; on Linux it signals the child when the parent exits.

	// Example: pdeathsig
	out, _ := execx.Command("printf", "ok").Pdeathsig(syscall.SIGTERM).Output()
	fmt.Print(out)
	// ok
}
