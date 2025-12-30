//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Terminate kills the process immediately.

	// Example: terminate
	proc := execx.Command("sleep", "2").Start()
	_ = proc.Terminate()
	res, _ := proc.Wait()
	fmt.Println(res.ExitCode != 0)
	// #bool true
}
