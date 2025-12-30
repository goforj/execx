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
	_, err := execx.Command("go", "env", "GOOS").
		StdoutWriter(&out).
		Run()
	fmt.Println(err == nil && out.Len() > 0)
	// #bool true
}
