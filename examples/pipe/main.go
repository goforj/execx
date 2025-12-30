//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"os"
	"strings"
)

func main() {
	// Pipe appends a new command to the pipeline.

	// Example: pipe
	if os.Getenv("EXECX_EXAMPLE_CHILD") == "1" {
		switch os.Getenv("EXECX_EXAMPLE_MODE") {
		case "emit":
			fmt.Print("go")
		case "upper":
			buf := make([]byte, 8)
			n, _ := os.Stdin.Read(buf)
			fmt.Print(strings.ToUpper(string(buf[:n])))
		}
		return
	}
	out, _ := execx.Command(os.Args[0]).
		Env("EXECX_EXAMPLE_CHILD=1", "EXECX_EXAMPLE_MODE=emit").
		Pipe(os.Args[0]).
		Env("EXECX_EXAMPLE_CHILD=1", "EXECX_EXAMPLE_MODE=upper").
		OutputTrimmed()
	fmt.Println(out)
	// #string GO
}
