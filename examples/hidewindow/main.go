//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// HideWindow is a no-op on non-Windows platforms.

	// Example: hide window
	fmt.Println(execx.Command("go", "env", "GOOS").HideWindow(true) != nil)
	// #bool true
	// Example: hide window
	fmt.Println(execx.Command("go", "env", "GOOS").HideWindow(true) != nil)
	// #bool true
}
