//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// StdinString sets stdin from a string.

	// Example: stdin string
	out, _ := execx.Command("cat").
		StdinString("hi").
		Output()
	fmt.Println(out)
	// #string hi
}
