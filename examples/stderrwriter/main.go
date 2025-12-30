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
	fmt.Println(err == nil && out.Len() > 0)
	// #bool true
}
