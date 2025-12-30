//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
	"time"
)

func main() {
	// GracefulShutdown sends a signal and escalates to kill after the timeout.

	// Example: graceful shutdown
	proc := execx.Command("sleep", "2").
		Start()
	_ = proc.GracefulShutdown(os.Interrupt, 100*time.Millisecond)
	res := proc.Wait()
	fmt.Println(res.ExitCode != 0)
	// #bool true
}
