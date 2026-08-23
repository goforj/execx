package execx

import "testing"

var benchmarkArgs []string

// benchmarkPayload prevents decoded benchmark results from being optimized away.
var benchmarkPayload testPayload

// BenchmarkCommandConstruction measures fluent argument assembly without subprocess noise.
func BenchmarkCommandConstruction(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkArgs = Command("tool", "base").
			Arg("value", []string{"one", "two"}, map[string]string{"--mode": "fast"}).
			Args()
	}
}

// BenchmarkShellEscaped measures the logging representation used by shadow printing.
func BenchmarkShellEscaped(b *testing.B) {
	cmd := Command("tool", "hello world", "it's", "$TOKEN")
	b.ReportAllocs()
	for b.Loop() {
		benchmarkArgs = []string{cmd.ShellEscaped()}
	}
}

// BenchmarkDecodeResult compares caller-owned and generic result decoding through the same command path.
func BenchmarkDecodeResult(b *testing.B) {
	b.Run("Into", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			var out testPayload
			if err := Command("printf", `{"name":"gopher"}`).DecodeJSON().Into(&out); err != nil {
				b.Fatalf("Into: %v", err)
			}
			benchmarkPayload = out
		}
	})
	b.Run("As", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			out, err := Command("printf", `{"name":"gopher"}`).DecodeJSON().As[testPayload]()
			if err != nil {
				b.Fatalf("As: %v", err)
			}
			benchmarkPayload = out
		}
	})
}
