//go:build linux

package execx

import "syscall"

// Setpgid sets the process group ID behavior.
func (c *Cmd) Setpgid(on bool) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.Setpgid = on
	return c
}

// Setsid sets the session ID behavior.
func (c *Cmd) Setsid(on bool) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.Setsid = on
	return c
}

// Pdeathsig sets a parent-death signal on Linux.
func (c *Cmd) Pdeathsig(sig syscall.Signal) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.Pdeathsig = sig
	return c
}

func (c *Cmd) ensureSysProcAttr() {
	if c.sysProcAttr == nil {
		c.sysProcAttr = &syscall.SysProcAttr{}
	}
}
