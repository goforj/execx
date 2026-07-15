//go:build plan9

package execx

import "syscall"

// Setpgid is a no-op because Plan 9 does not expose Unix process groups.
// @group OS Controls
func (c *Cmd) Setpgid(_ bool) *Cmd {
	return c
}

// Setsid is a no-op because Plan 9 does not expose Unix sessions.
// @group OS Controls
func (c *Cmd) Setsid(_ bool) *Cmd {
	return c
}

// Pdeathsig is a no-op because parent-death signals are Linux-specific.
// @group OS Controls
func (c *Cmd) Pdeathsig(_ syscall.Note) *Cmd {
	return c
}
