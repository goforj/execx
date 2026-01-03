//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// FromStderr decodes from stderr.

	// Example: decode from stderr
	type payload struct {
		Name string `json:"name"`
	}
	var out payload
	_ = execx.Command("sh", "-c", `printf '{"name":"gopher"}' 1>&2`).
		DecodeJSON().
		FromStderr().
		Into(&out)
	fmt.Println(out.Name)
	// #string gopher
}
