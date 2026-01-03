//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/goforj/execx"
	"strings"
)

func main() {
	// FromStderr decodes from stderr.

	// Example: decode from stderr
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
	_ = execx.Command("sh", "-c", "printf 'name=gopher' 1>&2").
		Decode(decoder).
		FromStderr().
		Into(&out)
	fmt.Println(out.Name)
	// #string gopher
}
