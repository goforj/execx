//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Setsid starts the child in a new session, detaching it from the terminal.

	// Example: setsid
	out, _ := execx.Command("printf", "ok").Setsid(true).Output()
	fmt.Print(out)
	// ok
}
