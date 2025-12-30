//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// StdoutWriter sets a raw writer for stdout.

	// Example: stdout writer
	var out strings.Builder
	_, _ = execx.Command("printf", "hello").
		StdoutWriter(&out).
		Run()
	fmt.Print(out.String())
	// hello
}
