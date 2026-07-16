package main

import (
	"fmt"
	"github.com/goforj/execx"
	"time"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// WithTimeout binds the command to a timeout. Replacing a timeout retains the
	// context previously supplied through WithContext as its parent.

	// Example: with timeout
	res, _ := execx.Command("go", "env", "GOOS").WithTimeout(2 * time.Second).Run()
	fmt.Println(res.ExitCode == 0)
	// #bool true
}
