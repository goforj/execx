//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// StdinBytes sets stdin from bytes.

	// Example: stdin bytes
	out, _ := execx.Command("cat").
		StdinBytes([]byte("hi")).
		Output()
	fmt.Println(out)
	// #string hi
}
