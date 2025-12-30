//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Setpgid sets the process group ID behavior.

	// Example: setpgid
	fmt.Println(execx.Command("go", "env", "GOOS").Setpgid(true) != nil)
	// #bool true
	// Example: setpgid
	fmt.Println(execx.Command("go", "env", "GOOS").Setpgid(true) != nil)
	// #bool true
	// Example: setpgid
	fmt.Println(execx.Command("go", "env", "GOOS").Setpgid(true) != nil)
	// #bool true
}
