//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// PipeStrict sets strict pipeline semantics.

	// Example: strict
	res, err := execx.Command("false").
		Pipe("printf", "ok").
		PipeStrict().
		Run()
	fmt.Println(err == nil && res.ExitCode != 0)
	// #bool true
}
