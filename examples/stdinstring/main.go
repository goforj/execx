package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// StdinString sets stdin from a string.

	// Example: stdin string
	out, _ := execx.Command("cat").
		StdinString("hi").
		Output()
	fmt.Println(out)
	// #string hi
}
