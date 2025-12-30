//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// PipeBestEffort sets best-effort pipeline semantics.

	// Example: best effort
	res, err := execx.Command("false").
		Pipe("printf", "ok").
		PipeBestEffort().
		Run()
	fmt.Println(err == nil && res.Stdout == "ok")
	// #bool true
}
