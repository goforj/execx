package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestHasExportedReceiver verifies that example generation follows the public API boundary.
func TestHasExportedReceiver(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "surface.go", `
package surface
type Public struct{}
type private struct{}
func Top() {}
func (*Public) Method() {}
func (*private) Hidden() {}
`, 0)
	if err != nil {
		t.Fatalf("parse declarations: %v", err)
	}

	got := map[string]bool{}
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			got[fn.Name.Name] = hasExportedReceiver(fn)
		}
	}
	if !got["Top"] || !got["Method"] || got["Hidden"] {
		t.Fatalf("unexpected receiver visibility: %v", got)
	}
}
