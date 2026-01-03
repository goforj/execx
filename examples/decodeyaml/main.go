//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// DecodeYAML configures YAML decoding for this command.
	// Decoding reads from stdout by default; use FromStdout, FromStderr, or FromCombined to select a source.

	// Example: decode yaml
	type payload struct {
		Name string `yaml:"name"`
	}
	var out payload
	_ = execx.Command("printf", "name: gopher").
		DecodeYAML().
		Into(&out)
	fmt.Println(out.Name)
	// #string gopher
}
