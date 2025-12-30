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
		OnStderr(func(line string) { lines = append(lines, line) }).
		Run()
	fmt.Println(err == nil && len(lines) == 1)
	// #bool true
}
