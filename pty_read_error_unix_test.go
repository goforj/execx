//go:build linux || darwin

package execx

import (
	"errors"
	"io"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

// ptyErrorReader provides deterministic terminal read failures for the platform adapter tests.
type ptyErrorReader struct {
	err error
}

// Read returns the configured error so PTY EOF handling can be tested without kernel timing.
func (r ptyErrorReader) Read([]byte) (int, error) {
	return 0, r.err
}

// ptyDataErrorReader returns output and an error in the same read.
type ptyDataErrorReader struct {
	data string
	err  error
}

// Read reproduces a reader that makes progress while reporting a transient condition.
func (r ptyDataErrorReader) Read(p []byte) (int, error) {
	return copy(p, r.data), r.err
}

// ptyRetryReader returns a transient failure before making output available.
type ptyRetryReader struct {
	attempts int
}

// Read reproduces the macOS PTY read sequence without requiring kernel timing.
func (r *ptyRetryReader) Read(p []byte) (int, error) {
	r.attempts++
	if r.attempts == 1 {
		return 0, syscall.EAGAIN
	}
	return copy(p, "ready"), nil
}

// ptyExitDrainReader makes final output visible only after the first unavailable read.
type ptyExitDrainReader struct {
	attempts int
}

// Read reproduces final PTY output racing with child process exit.
func (r *ptyExitDrainReader) Read(p []byte) (int, error) {
	r.attempts++
	if r.attempts == 2 {
		return copy(p, "final"), nil
	}
	return 0, syscall.EAGAIN
}

// testPTYOutputReader supplies deterministic readiness results to the Unix PTY adapter.
func testPTYOutputReader(reader io.Reader, hungUp bool, waitErr error) io.Reader {
	return ptyEOFReader{
		reader: reader,
		waitReadable: func() (bool, error) {
			return hungUp, waitErr
		},
	}
}

// TestPTYEOFReader distinguishes a normal Unix terminal hangup from a real read failure.
func TestPTYEOFReader(t *testing.T) {
	reader := testPTYOutputReader(ptyErrorReader{err: syscall.EIO}, false, nil)
	if _, err := reader.Read(make([]byte, 1)); !errors.Is(err, io.EOF) {
		t.Fatalf("ptyOutputReader(EIO) error = %v, want EOF", err)
	}

	want := errors.New("read failed")
	reader = testPTYOutputReader(ptyErrorReader{err: want}, false, nil)
	if _, err := reader.Read(make([]byte, 1)); !errors.Is(err, want) {
		t.Fatalf("ptyOutputReader() error = %v, want %v", err, want)
	}
}

// TestPTYEOFReaderRetriesTemporaryUnavailability keeps an interactive PTY alive across transient reads.
func TestPTYEOFReaderRetriesTemporaryUnavailability(t *testing.T) {
	source := &ptyRetryReader{}
	reader := testPTYOutputReader(source, false, nil)
	buf := make([]byte, len("ready"))

	n, err := reader.Read(buf)
	if err != nil {
		t.Fatalf("ptyOutputReader(EAGAIN) error = %v", err)
	}
	if got := string(buf[:n]); got != "ready" {
		t.Fatalf("ptyOutputReader(EAGAIN) output = %q, want %q", got, "ready")
	}
	if source.attempts != 2 {
		t.Fatalf("ptyOutputReader(EAGAIN) attempts = %d, want 2", source.attempts)
	}
}

// TestPTYEOFReaderPreservesDataWithTemporaryUnavailability keeps readable output from being discarded.
func TestPTYEOFReaderPreservesDataWithTemporaryUnavailability(t *testing.T) {
	reader := testPTYOutputReader(
		ptyDataErrorReader{data: "ready", err: syscall.EAGAIN},
		false,
		nil,
	)
	buf := make([]byte, len("ready"))

	n, err := reader.Read(buf)
	if err != nil {
		t.Fatalf("ptyOutputReader(data with EAGAIN) error = %v", err)
	}
	if got := string(buf[:n]); got != "ready" {
		t.Fatalf("ptyOutputReader(data with EAGAIN) output = %q, want %q", got, "ready")
	}
}

// TestPTYEOFReaderDrainsOutputAfterHangup preserves buffered output before reporting terminal EOF.
func TestPTYEOFReaderDrainsOutputAfterHangup(t *testing.T) {
	source := &ptyExitDrainReader{}
	reader := testPTYOutputReader(source, true, nil)
	buf := make([]byte, len("final"))

	n, err := reader.Read(buf)
	if err != nil {
		t.Fatalf("ptyOutputReader(final output after hangup) error = %v", err)
	}
	if got := string(buf[:n]); got != "final" {
		t.Fatalf("ptyOutputReader(final output after hangup) output = %q, want %q", got, "final")
	}
	if _, err := reader.Read(buf); !errors.Is(err, io.EOF) {
		t.Fatalf("ptyOutputReader(EAGAIN after hangup drain) error = %v, want EOF", err)
	}
}

// TestPTYEOFReaderPreservesReadinessFailure keeps polling failures visible to callers.
func TestPTYEOFReaderPreservesReadinessFailure(t *testing.T) {
	want := errors.New("poll failed")
	reader := testPTYOutputReader(ptyErrorReader{err: syscall.EAGAIN}, false, want)

	if _, err := reader.Read(make([]byte, 1)); !errors.Is(err, want) {
		t.Fatalf("ptyOutputReader(EAGAIN poll) error = %v, want %v", err, want)
	}
}

// TestWaitPTYReadableDistinguishesPTYDataFromHangup exercises the platform PTY polling contract.
func TestWaitPTYReadableDistinguishesPTYDataFromHangup(t *testing.T) {
	master, slave, err := openPTY()
	if err != nil {
		t.Fatalf("openPTY() error = %v", err)
	}
	defer master.Close()
	defer slave.Close()
	reader, ok := ptyOutputReader(master).(ptyEOFReader)
	if !ok {
		t.Fatalf("ptyOutputReader() type = %T, want ptyEOFReader", ptyOutputReader(master))
	}

	if _, err := slave.WriteString("ready"); err != nil {
		t.Fatalf("slave.WriteString() error = %v", err)
	}
	hungUp, err := reader.waitReadable()
	if err != nil {
		t.Fatalf("waitPTYReadable(data) error = %v", err)
	}
	if hungUp {
		t.Fatal("waitPTYReadable(data) hungUp = true, want false")
	}
	buf := make([]byte, len("ready"))
	if _, err := io.ReadFull(master, buf); err != nil {
		t.Fatalf("io.ReadFull(master) error = %v", err)
	}
	if got := string(buf); got != "ready" {
		t.Fatalf("master output = %q, want %q", got, "ready")
	}

	if err := slave.Close(); err != nil {
		t.Fatalf("slave.Close() error = %v", err)
	}
	hungUp, err = reader.waitReadable()
	if err != nil {
		t.Fatalf("waitPTYReadable(hangup) error = %v", err)
	}
	if !hungUp {
		t.Fatal("waitPTYReadable(hangup) hungUp = false, want true")
	}
}

// TestWaitPTYReadableRetriesInterruptedPoll preserves blocking semantics across Unix signals.
func TestWaitPTYReadableRetriesInterruptedPoll(t *testing.T) {
	master, slave, err := openPTY()
	if err != nil {
		t.Fatalf("openPTY() error = %v", err)
	}
	defer master.Close()
	defer slave.Close()

	attempts := 0
	hungUp, err := waitPTYReadableWith(master, func(descriptors []unix.PollFd, timeout int) (int, error) {
		attempts++
		if attempts == 1 {
			return 0, syscall.EINTR
		}
		descriptors[0].Revents = unix.POLLHUP
		return 1, nil
	})
	if err != nil {
		t.Fatalf("waitPTYReadableWith(EINTR) error = %v", err)
	}
	if !hungUp {
		t.Fatal("waitPTYReadableWith(EINTR) hungUp = false, want true")
	}
	if attempts != 2 {
		t.Fatalf("waitPTYReadableWith(EINTR) attempts = %d, want 2", attempts)
	}
}

// TestWaitPTYReadablePreservesPollFailure keeps kernel polling errors visible.
func TestWaitPTYReadablePreservesPollFailure(t *testing.T) {
	master, slave, err := openPTY()
	if err != nil {
		t.Fatalf("openPTY() error = %v", err)
	}
	defer master.Close()
	defer slave.Close()

	want := errors.New("poll failed")
	_, err = waitPTYReadableWith(master, func([]unix.PollFd, int) (int, error) {
		return 0, want
	})
	if !errors.Is(err, want) {
		t.Fatalf("waitPTYReadableWith() error = %v, want %v", err, want)
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
