//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Setsid sets the session ID behavior.

	// Example: setsid
	fmt.Println(execx.Command("go", "env", "GOOS").Setsid(true) != nil)
	// #bool true
	// Example: setsid
	fmt.Println(execx.Command("go", "env", "GOOS").Setsid(true) != nil)
	// #bool true
	// Example: setsid
	fmt.Println(execx.Command("go", "env", "GOOS").Setsid(true) != nil)
	// #bool true
}
