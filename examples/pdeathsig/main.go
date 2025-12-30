//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Pdeathsig sets a parent-death signal on Linux.

	// Example: pdeathsig
	fmt.Println(execx.Command("go", "env", "GOOS").Pdeathsig(0) != nil)
	// #bool true
	// Example: pdeathsig
	fmt.Println(execx.Command("go", "env", "GOOS").Pdeathsig(0) != nil)
	// #bool true
	// Example: pdeathsig
	fmt.Println(execx.Command("go", "env", "GOOS").Pdeathsig(0) != nil)
	// #bool true
}
