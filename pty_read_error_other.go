//go:build !linux && !darwin

package execx

import "io"

// ptyOutputReader preserves native read failures on platforms without Unix PTY hangup semantics.
func ptyOutputReader(reader io.Reader) io.Reader {
	return reader
}
