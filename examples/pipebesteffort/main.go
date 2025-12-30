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
	// PipeBestEffort sets best-effort pipeline semantics.

	// Example: best effort
	if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
		switch os.Getenv("EXECX_EXAMPLE_MODE") {
		case "sleep":
			time.Sleep(200 * time.Millisecond)
		case "ok":
			fmt.Print("ok")
		}
		return
	}
	res := execx.Command(os.Args[0]).
		Env("EXECX_EXAMPLE_CHILD=1", "EXECX_EXAMPLE_MODE=sleep").
		WithTimeout(50 * time.Millisecond).
		Pipe(os.Args[0]).
		Env("EXECX_EXAMPLE_CHILD=1", "EXECX_EXAMPLE_MODE=ok").
		PipeBestEffort().
		Run()
	fmt.Println(res.Stdout)
	// #string ok
}
