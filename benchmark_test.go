package execx

import "testing"

var benchmarkArgs []string

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
