//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Setsid is a no-op on Windows; on Unix it starts a new session.

	// Example: setsid
	out, _ := execx.Command("printf", "ok").Setsid(true).Output()
	fmt.Print(out)
	// ok
}
