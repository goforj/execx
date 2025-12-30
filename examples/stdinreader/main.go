//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// StdinReader sets stdin from an io.Reader.

	// Example: stdin reader
	out, _ := execx.Command("cat").
		StdinReader(strings.NewReader("hi")).
		Output()
	fmt.Println(out)
	// #string hi
}
