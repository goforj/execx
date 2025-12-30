//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
)

func main() {
	// Interrupt sends an interrupt signal to the process.

	// Example: interrupt
	proc := execx.Command("sleep", "2").Start()
	_ = proc.Interrupt()
	res, _ := proc.Wait()
	fmt.Println(res.IsSignal(os.Interrupt))
	// #bool true
}
