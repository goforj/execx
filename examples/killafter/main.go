//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"time"
)

func main() {
	// KillAfter terminates the process after the given duration.

	// Example: kill after
	proc := execx.Command("sleep", "2").Start()
	proc.KillAfter(100 * time.Millisecond)
	res, _ := proc.Wait()
	fmt.Println(res.ExitCode != 0)
	// #bool true
}
