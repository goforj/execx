//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Setpgid is a no-op on Windows.

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
