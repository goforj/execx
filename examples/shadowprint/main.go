//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// ShadowPrint configures shadow printing for this command chain.

	// Example: shadow print
	_, _ = execx.Command("echo", "hello world").
		ShadowPrint().
		OnStdout(func(line string) { fmt.Println(line) }).
		Run()
	// execx > bash -c 'echo "hello world"'
	// hello world
	// execx > bash -c 'echo "hello world"' (1ms)

	// Example: shadow print options
	mask := func(cmd string) string {
		return strings.ReplaceAll(cmd, "token", "***")
	}
	formatter := func(ev execx.ShadowEvent) string {
		return fmt.Sprintf("shadow: %s %s", ev.Phase, ev.Command)
	}
	_, _ = execx.Command("echo", "hello world").
		ShadowPrint(
			execx.WithPrefix("execx"),
			execx.WithMask(mask),
			execx.WithFormatter(formatter),
		).
		OnStdout(func(line string) { fmt.Println(line) }).
		Run()
	// shadow: before bash -c 'echo "hello world"'
	// hello world
	// shadow: after bash -c 'echo "hello world"'
}
