package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Error returns the wrapped error message when available.

	// Example: error string
	err := execx.ErrExec{Err: fmt.Errorf("boom")}
	fmt.Println(err.Error())
	// #string boom
}
