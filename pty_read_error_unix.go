//go:build linux || darwin

package execx

import (
	"errors"
	"io"
	"syscall"
)

// ptyEOFReader converts the Unix terminal hangup convention without changing other read errors.
type ptyEOFReader struct {
	reader io.Reader
}

// Read treats a Unix PTY hangup as EOF because the slave has already closed normally.
func (r ptyEOFReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if errors.Is(err, syscall.EIO) {
		return n, io.EOF
	}
	return n, err
}

// ptyOutputReader limits terminal hangup normalization to reads so writer failures remain visible.
func ptyOutputReader(reader io.Reader) io.Reader {
	return ptyEOFReader{reader: reader}
}
