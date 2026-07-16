package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Setpgid places the child in a new process group for group signals.

	// Example: setpgid
	out, _ := execx.Command("printf", "ok").Setpgid(true).Output()
	fmt.Print(out)
	// ok
}
