//go:build linux

package execx

import "syscall"

// Setpgid places the child in a new process group for group signals.
// @group OS Controls
//
// Example: setpgid
//
//	out, _ := execx.Command("printf", "ok").Setpgid(true).Output()
//	fmt.Print(out)
//	// ok
func (c *Cmd) Setpgid(on bool) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.Setpgid = on
	return c
}

// Setsid starts the child in a new session, detaching it from the terminal.
// @group OS Controls
//
// Example: setsid
//
//	out, _ := execx.Command("printf", "ok").Setsid(true).Output()
//	fmt.Print(out)
//	// ok
func (c *Cmd) Setsid(on bool) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.Setsid = on
	return c
}

// Pdeathsig sets a parent-death signal on Linux so the child is signaled if the parent exits.
// @group OS Controls
//
// Example: pdeathsig
//
//	out, _ := execx.Command("printf", "ok").Pdeathsig(syscall.SIGTERM).Output()
//	fmt.Print(out)
//	// ok
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
