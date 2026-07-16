package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// StdinFile sets stdin from a file.

	// Example: stdin file
	file, _ := os.CreateTemp("", "execx-stdin")
	_, _ = file.WriteString("hi")
	_, _ = file.Seek(0, 0)
	out, _ := execx.Command("cat").
		StdinFile(file).
		Output()
	fmt.Println(out)
	// #string hi
}
