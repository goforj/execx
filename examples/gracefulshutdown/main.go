package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
	"time"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// GracefulShutdown sends a signal and escalates to kill after the timeout.

	// Example: graceful shutdown
	proc := execx.Command("sleep", "2").Start()
	_ = proc.GracefulShutdown(os.Interrupt, 100*time.Millisecond)
	res, _ := proc.Wait()
	fmt.Println(res.IsSignal(os.Interrupt))
	// #bool true
}
