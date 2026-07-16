package main

import (
	"context"
	"fmt"
	"github.com/goforj/execx"
	"time"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// WithContext binds the command to a context. A nil context is normalized to
	// context.Background when execution begins.

	// Example: with context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, _ := execx.Command("go", "env", "GOOS").WithContext(ctx).Run()
	fmt.Println(res.ExitCode == 0)
	// #bool true
}
