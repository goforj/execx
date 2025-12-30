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
	proc := execx.Command("sleep", "2").Start()
	_ = proc.Send(os.Interrupt)
	res, _ := proc.Wait()
	fmt.Println(res.IsSignal(os.Interrupt))
	// #bool true
}
