package main

import (
	"fmt"
	"github.com/goforj/execx"
	"time"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// WithDeadline binds the command to a deadline. Replacing a deadline retains the
	// context previously supplied through WithContext as its parent.

	// Example: with deadline
	res, _ := execx.Command("go", "env", "GOOS").WithDeadline(time.Now().Add(2 * time.Second)).Run()
	fmt.Println(res.ExitCode == 0)
	// #bool true
}
