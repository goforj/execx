package main

import (
	"fmt"
	"github.com/goforj/execx"
)

// main keeps this documented example executable so API drift fails during compilation.
func main() {
	// As executes the command and decodes its selected output into T.

	// Example: decode as a value
	type payload struct {
		Name string `json:"name"`
	}
	out, err := execx.Command("printf", `{"name":"gopher"}`).
		DecodeJSON().
		As[payload]()
	fmt.Println(err == nil, out.Name)
	// true gopher
}
