package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// FromCombined decodes from combined stdout+stderr.

	// Example: decode combined
	type payload struct {
		Name string `json:"name"`
	}
	var out payload
	_ = execx.Command("sh", "-c", `printf '{"name":"gopher"}'`).
		DecodeJSON().
		FromCombined().
		Into(&out)
	fmt.Println(out.Name)
	// #string gopher
}
