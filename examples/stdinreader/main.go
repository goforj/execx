package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// StdinReader sets stdin from an io.Reader.

	// Example: stdin reader
	out, _ := execx.Command("cat").
		StdinReader(strings.NewReader("hi")).
		Output()
	fmt.Println(out)
	// #string hi
}
