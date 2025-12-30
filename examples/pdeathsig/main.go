//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Pdeathsig sets a parent-death signal on Linux so the child is signaled if the parent exits.

	// Example: pdeathsig
	out, _ := execx.Command("printf", "ok").Pdeathsig(syscall.SIGTERM).Output()
	fmt.Print(out)
	// ok
}
