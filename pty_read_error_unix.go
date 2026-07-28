//go:build linux || darwin

package execx

import (
	"errors"
	"io"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

// ptyEOFReader converts the Unix terminal hangup convention without changing other read errors.
type ptyEOFReader struct {
	reader       io.Reader
	waitReadable func() (bool, error)
}

// Read treats terminal hangups as EOF and preserves blocking-reader semantics for transient unavailability.
func (r ptyEOFReader) Read(p []byte) (int, error) {
	hungUp := false
	for {
		n, err := r.reader.Read(p)
		if errors.Is(err, syscall.EIO) {
			return n, io.EOF
		}
		if errors.Is(err, syscall.EAGAIN) {
			if n > 0 {
				return n, nil
			}
			if hungUp {
				return 0, io.EOF
			}
			hungUp, err = r.waitReadable()
			if err != nil {
				return 0, err
			}
			continue
		}
		return n, err
	}
}

// ptyOutputReader limits terminal hangup normalization to reads so writer failures remain visible.
func ptyOutputReader(file *os.File) io.Reader {
	return ptyEOFReader{
		reader: file,
		waitReadable: func() (bool, error) {
			return waitPTYReadable(file)
		},
	}
}

// waitPTYReadable blocks until the PTY has output or all slave descriptors have closed.
func waitPTYReadable(file *os.File) (bool, error) {
	return waitPTYReadableWith(file, unix.Poll)
}

// waitPTYReadableWith delegates polling so interrupted and failed waits remain deterministic in tests.
func waitPTYReadableWith(file *os.File, poll func([]unix.PollFd, int) (int, error)) (bool, error) {
	descriptors := []unix.PollFd{{
		Fd:     int32(file.Fd()),
		Events: unix.POLLIN | unix.POLLHUP,
	}}
	for {
		_, err := poll(descriptors, -1)
		if errors.Is(err, syscall.EINTR) {
			continue
		}
		if err != nil {
			return false, err
		}
		return descriptors[0].Revents&unix.POLLHUP != 0, nil
	}
}
