//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// FromStdout decodes from stdout (default).

	// Example: decode from stdout
	type payload struct {
		Name string `json:"name"`
	}
	var out payload
	_ = execx.Command("printf", `{"name":"gopher"}`).
		DecodeJSON().
		FromStdout().
		Into(&out)
	fmt.Println(out.Name)
	// #string gopher
}
