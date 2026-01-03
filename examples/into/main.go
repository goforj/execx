//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Into executes the command and decodes into dst.

	// Example: decode into
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
