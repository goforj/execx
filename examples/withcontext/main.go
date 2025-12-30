//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"github.com/goforj/execx"
	"time"
)

func main() {
	// WithContext binds the command to a context.

	// Example: with context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := execx.Command("go", "env", "GOOS").WithContext(ctx).Run()
	fmt.Println(err == nil && res.ExitCode == 0)
	// #bool true
}
