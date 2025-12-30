//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"github.com/goforj/execx"
	"log"
	"time"
)

func main() {
	// Run executes the command and returns the result and any error.

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := execx.
		Command("printf", "hello\nworld\n").
		Pipe("tr", "a-z", "A-Z").
		Env("MODE=demo").
		WithContext(ctx).
		OnStdout(func(line string) {
			fmt.Println("OUT:", line)
		}).
		OnStderr(func(line string) {
			fmt.Println("ERR:", line)
		}).
		Run()

	if !res.OK() {
		log.Fatalf("command failed: %v", err)
	}

	fmt.Printf("Stdout: %q\n", res.Stdout)
	fmt.Printf("Stderr: %q\n", res.Stderr)
	fmt.Printf("ExitCode: %d\n", res.ExitCode)
	fmt.Printf("Error: %v\n", res.Err)
	fmt.Printf("Duration: %v\n", res.Duration)
	// OUT: HELLO
	// OUT: WORLD
	// Stdout: "HELLO\nWORLD\n"
	// Stderr: ""
	// ExitCode: 0
	// Error: <nil>
	// Duration: 10.123456ms

}
