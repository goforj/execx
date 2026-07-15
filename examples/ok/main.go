package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// OK reports whether the command exited cleanly without errors.

	// Example: ok
	res, _ := execx.Command("go", "env", "GOOS").Run()
	fmt.Println(res.OK())
	// #bool true
}
