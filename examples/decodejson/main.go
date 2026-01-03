//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// DecodeJSON configures JSON decoding for this command.
	// Decoding reads from stdout by default; use FromStdout, FromStderr, or FromCombined to select a source.

	// Example: decode json
	type payload struct {
		Name string `json:"name"`
	}
	var out payload
	_ = execx.Command("printf", `{"name":"gopher"}`).
		DecodeJSON().
		Into(&out)
	fmt.Println(out.Name)
	// #string gopher
}
