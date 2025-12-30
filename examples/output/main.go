//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Output executes the command and returns stdout and any error.

	// Example: output
	out, _ := execx.Command("go", "env", "GOOS").Output()
	fmt.Println(out != "")
	// #bool true
}
