//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"time"
)

func main() {
	// WithTimeout binds the command to a timeout.

	// Example: with timeout
	res, err := execx.Command("go", "env", "GOOS").WithTimeout(2 * time.Second).Run()
	fmt.Println(err == nil && res.ExitCode == 0)
	// #bool true
}
