//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// Trim trims whitespace before decoding.

	// Example: decode trim
	type payload struct {
		Name string `json:"name"`
	}
	var out payload
	_ = execx.Command("printf", "  {\"name\":\"gopher\"}  ").
		DecodeJSON().
		Trim().
		Into(&out)
	fmt.Println(out.Name)
	// #string gopher
}
