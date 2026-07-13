package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// PipeStrict sets strict pipeline semantics (stop on first failure).

	// Example: strict
	res, _ := execx.Command("false").
		Pipe("printf", "ok").
		PipeStrict().
		Run()
	fmt.Println(res.ExitCode != 0)
	// #bool true
}
