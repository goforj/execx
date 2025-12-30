//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// OK reports whether the command exited cleanly without errors.

	// Example: ok
	res, err := execx.Command("go", "env", "GOOS").Run()
	fmt.Println(err == nil && res.OK())
	// #bool true
}
