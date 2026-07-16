package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// StdinBytes sets stdin from a copy of bytes so later caller mutation cannot
	// change the command input.

	// Example: stdin bytes
	out, _ := execx.Command("cat").
		StdinBytes([]byte("hi")).
		Output()
	fmt.Println(out)
	// #string hi
}
