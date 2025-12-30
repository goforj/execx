//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// PipeBestEffort sets best-effort pipeline semantics (run all stages, surface the first error).

	// Example: best effort
	res, _ := execx.Command("false").
		Pipe("printf", "ok").
		PipeBestEffort().
		Run()
	fmt.Print(res.Stdout)
	// ok
}
