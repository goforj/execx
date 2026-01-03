//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// Decode configures a custom decoder for this command.

	// Example: decode custom
	type payload struct {
		Name string
	}
	decoder := execx.DecoderFunc(func(data []byte, dst any) error {
		out, ok := dst.(*payload)
		if !ok {
			return fmt.Errorf("expected *payload")
		}
		_, val, ok := strings.Cut(string(data), "=")
		if !ok {
			return fmt.Errorf("invalid payload")
		}
		out.Name = val
		return nil
	})
	var out payload
	_ = execx.Command("printf", "name=gopher").
		Decode(decoder).
		Into(&out)
	fmt.Println(out.Name)
	// #string gopher
}
