//go:build !unix && !windows && !plan9

package execx

import "syscall"

// Setpgid is a no-op on platforms without Unix process groups.
// @group OS Controls
func (c *Cmd) Setpgid(_ bool) *Cmd {
	return c
}

// Setsid is a no-op on platforms without Unix sessions.
// @group OS Controls
func (c *Cmd) Setsid(_ bool) *Cmd {
	return c
}

// Pdeathsig is a no-op outside Linux because parent-death signals are Linux-specific.
// @group OS Controls
func (c *Cmd) Pdeathsig(_ syscall.Signal) *Cmd {
	return c
}
