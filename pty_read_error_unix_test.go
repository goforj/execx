//go:build linux || darwin

package execx

import (
	"errors"
	"io"
	"syscall"
	"testing"
)

// ptyErrorReader provides deterministic terminal read failures for the platform adapter tests.
type ptyErrorReader struct {
	err error
}

// Read returns the configured error so PTY EOF handling can be tested without kernel timing.
func (r ptyErrorReader) Read([]byte) (int, error) {
	return 0, r.err
}

// TestPTYEOFReader distinguishes a normal Unix terminal hangup from a real read failure.
func TestPTYEOFReader(t *testing.T) {
	reader := ptyOutputReader(ptyErrorReader{err: syscall.EIO})
	if _, err := reader.Read(make([]byte, 1)); !errors.Is(err, io.EOF) {
		t.Fatalf("ptyOutputReader(EIO) error = %v, want EOF", err)
	}

	want := errors.New("read failed")
	reader = ptyOutputReader(ptyErrorReader{err: want})
	if _, err := reader.Read(make([]byte, 1)); !errors.Is(err, want) {
		t.Fatalf("ptyOutputReader() error = %v, want %v", err, want)
	}
}

// TestWithPTYTreatsHangupAsEOF exercises the kernel PTY path that pipes cannot reproduce.
func TestWithPTYTreatsHangupAsEOF(t *testing.T) {
	result, err := Command("printf", "hello").WithPTY().Run()
	if err != nil {
		t.Fatalf("WithPTY().Run() error = %v", err)
	}
	if result.Stdout != "hello" {
		t.Fatalf("WithPTY().Run() stdout = %q, want %q", result.Stdout, "hello")
	}
}

// TestWithPTYPreservesExitStatusAfterHangup keeps terminal EOF handling from hiding child failures.
func TestWithPTYPreservesExitStatusAfterHangup(t *testing.T) {
	result, err := Command("sh", "-c", "printf failed; exit 7").WithPTY().Run()
	if err != nil {
		t.Fatalf("WithPTY().Run() error = %v", err)
	}
	if result.ExitCode != 7 {
		t.Fatalf("WithPTY().Run() exit code = %d, want 7", result.ExitCode)
	}
	if result.Stdout != "failed" {
		t.Fatalf("WithPTY().Run() stdout = %q, want %q", result.Stdout, "failed")
	}
}
