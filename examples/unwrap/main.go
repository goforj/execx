package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Unwrap exposes the underlying error.

	// Example: unwrap
	err := execx.ErrExec{Err: fmt.Errorf("boom")}
	fmt.Println(err.Unwrap() != nil)
	// #bool true
}
