//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// DecodeJSON configures JSON decoding for this command.

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
