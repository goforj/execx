//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
)

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
