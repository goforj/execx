//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// CreationFlags sets Windows creation flags.

	// Example: creation flags
	fmt.Println(execx.Command("go", "env", "GOOS").CreationFlags(0) != nil)
	// #bool true
	// Example: creation flags
	fmt.Println(execx.Command("go", "env", "GOOS").CreationFlags(0) != nil)
	// #bool true
}
