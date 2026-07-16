package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Dir sets the working directory.

	// Example: change dir
	dir := os.TempDir()
	out, _ := execx.Command("pwd").
		Dir(dir).
		OutputTrimmed()
	fmt.Println(out == dir)
	// #bool true
}
