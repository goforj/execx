//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Interrupt sends an interrupt signal to the process.

	// Example: interrupt
	proc := execx.Command("sleep", "2").Start()
	_ = proc.Interrupt()
	res, err := proc.Wait()
	fmt.Println(err != nil || res.ExitCode != 0)
	// #bool true
}
