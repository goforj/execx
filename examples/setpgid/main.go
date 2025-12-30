//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Setpgid places the child in a new process group for group signals.

	// Example: setpgid
	out, _ := execx.Command("printf", "ok").Setpgid(true).Output()
	fmt.Print(out)
	// ok
}
