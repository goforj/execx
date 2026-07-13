package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// StdoutWriter sets a raw writer for stdout.
	//
	// When the writer is a terminal and no line callbacks or combined output are enabled,
	// execx passes stdout through directly and does not buffer it for results.

	// Example: stdout writer
	var out strings.Builder
	_, _ = execx.Command("printf", "hello").
		StdoutWriter(&out).
		Run()
	fmt.Print(out.String())
	// hello
}
