//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
)

func main() {
	// Send sends a signal to the process.

	// Example: send signal
	proc := execx.Command("sleep", "2").
		Start()
	_ = proc.Send(os.Interrupt)
	res, err := proc.Wait()
	fmt.Println(err != nil || res.ExitCode != 0)
	// #bool true
}
