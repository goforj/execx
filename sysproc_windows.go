//go:build windows

package execx

import "syscall"

// Setpgid is a no-op on Windows.
func (c *Cmd) Setpgid(_ bool) *Cmd {
	return c
}

// Setsid is a no-op on Windows.
func (c *Cmd) Setsid(_ bool) *Cmd {
	return c
}

// Pdeathsig is a no-op on Windows.
func (c *Cmd) Pdeathsig(_ syscall.Signal) *Cmd {
	return c
}

// CreationFlags sets Windows creation flags.
func (c *Cmd) CreationFlags(flags uint32) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.CreationFlags = flags
	return c
}

// HideWindow controls window visibility.
func (c *Cmd) HideWindow(on bool) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.HideWindow = on
	return c
}

func (c *Cmd) ensureSysProcAttr() {
	if c.sysProcAttr == nil {
		c.sysProcAttr = &syscall.SysProcAttr{}
	}
}
