package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// StderrWriter sets a raw writer for stderr.
	//
	// When the writer is a terminal and no line callbacks or combined output are enabled,
	// execx passes stderr through directly and does not buffer it for results.

	// Example: stderr writer
	var out strings.Builder
	_, err := execx.Command("go", "env", "-badflag").
		StderrWriter(&out).
		Run()
	fmt.Print(out.String())
	fmt.Println(err == nil)
	// flag provided but not defined: -badflag
	// usage: go env [-json] [-changed] [-u] [-w] [var ...]
	// Run 'go help env' for details.
	// false
}
