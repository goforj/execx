//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// OutputTrimmed executes the command and returns trimmed stdout and any error.

	// Example: output trimmed
	out, _ := execx.Command("go", "env", "GOOS").OutputTrimmed()
	fmt.Println(out != "")
	// #bool true
}
