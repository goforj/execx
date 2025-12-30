//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
)

func main() {
	// PipeStrict sets strict pipeline semantics.

	// Example: strict
	if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
		switch os.Getenv("EXECX_EXAMPLE_MODE") {
		case "fail":
			os.Exit(2)
		case "ok":
			fmt.Print("ok")
		}
		return
	}
	res := execx.Command(os.Args[0]).
		Env("EXECX_EXAMPLE_CHILD=1", "EXECX_EXAMPLE_MODE=fail").
		Pipe(os.Args[0]).
		Env("EXECX_EXAMPLE_CHILD=1", "EXECX_EXAMPLE_MODE=ok").
		PipeStrict().
		Run()
	fmt.Println(res.ExitCode)
	// #int 2
}
