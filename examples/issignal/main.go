//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
)

func main() {
	// IsSignal reports whether the command terminated due to a signal.

	// Example: signal
	res, err := execx.Command("go", "env", "GOOS").Run()
	fmt.Println(err == nil && res.IsSignal(os.Interrupt))
	// #bool false
}
