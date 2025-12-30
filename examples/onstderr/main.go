//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// OnStderr registers a line callback for stderr.

	// Example: stderr lines
	var lines []string
	_, err := execx.Command("go", "env", "-badflag").
		OnStderr(func(line string) {
			lines = append(lines, line)
			fmt.Println(line)
		}).
		Run()
	fmt.Println(err == nil)
	// flag provided but not defined: -badflag
	// usage: go env [-json] [-changed] [-u] [-w] [var ...]
	// Run 'go help env' for details.
	// false
}
