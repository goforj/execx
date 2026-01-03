//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
)

func main() {
	// DecodeYAML configures YAML decoding for this command.

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
