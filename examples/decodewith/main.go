//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// DecodeWith executes the command and decodes stdout into dst.

	// Example: decode with
	type payload struct {
		Name string `json:"name"`
	}
	var out payload
	_ = execx.Command("printf", `{"name":"gopher"}`).
		DecodeWith(&out, execx.DecoderFunc(json.Unmarshal))
	fmt.Println(out.Name)
	// #string gopher
}
