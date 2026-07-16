package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// Unwrap exposes the underlying error.

	// Example: unwrap
	err := execx.ErrExec{Err: fmt.Errorf("boom")}
	fmt.Println(err.Unwrap() != nil)
	// #bool true
}
