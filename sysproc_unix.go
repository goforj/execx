//go:build unix && !linux

package execx

import "syscall"

// Setpgid places the child in a new process group for group signals.
// @group OS Controls
func (c *Cmd) Setpgid(on bool) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.Setpgid = on
	return c
}

// Setsid starts the child in a new session, detaching it from the terminal.
// @group OS Controls
func (c *Cmd) Setsid(on bool) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.Setsid = on
	return c
}

// Pdeathsig is a no-op on non-Linux platforms; on Linux it signals the child when the parent exits.
// @group OS Controls
func (c *Cmd) Pdeathsig(_ syscall.Signal) *Cmd {
	return c
}

func (c *Cmd) ensureSysProcAttr() {
	if c.sysProcAttr == nil {
		c.sysProcAttr = &syscall.SysProcAttr{}
	}
}
