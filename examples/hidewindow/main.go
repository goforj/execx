//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// HideWindow controls window visibility and sets CREATE_NO_WINDOW for console apps.

	// Example: hide window
	fmt.Println(execx.Command("go", "env", "GOOS").HideWindow(true) != nil)
	// #bool true
	// Example: hide window
	fmt.Println(execx.Command("go", "env", "GOOS").HideWindow(true) != nil)
	// #bool true
}
