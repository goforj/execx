//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// StderrWriter sets a raw writer for stderr.

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
