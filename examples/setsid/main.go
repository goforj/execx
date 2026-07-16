package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Setsid starts the child in a new session, detaching it from the terminal.

	// Example: setsid
	out, _ := execx.Command("printf", "ok").Setsid(true).Output()
	fmt.Print(out)
	// ok
}
