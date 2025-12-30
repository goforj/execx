//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// CreationFlags is a no-op on non-Windows platforms.

	// Example: creation flags
	fmt.Println(execx.Command("go", "env", "GOOS").CreationFlags(0) != nil)
	// #bool true
	// Example: creation flags
	fmt.Println(execx.Command("go", "env", "GOOS").CreationFlags(0) != nil)
	// #bool true
}
