//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// OutputBytes executes the command and returns stdout bytes and any error.

	// Example: output bytes
	out, _ := execx.Command("go", "env", "GOOS").OutputBytes()
	fmt.Println(len(out) > 0)
	// #bool true
}
