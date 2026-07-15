package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// IsSignal reports whether the command terminated due to a signal.

	// Example: signal
	res, _ := execx.Command("go", "env", "GOOS").Run()
	fmt.Println(res.IsSignal(os.Interrupt))
	// false
}
