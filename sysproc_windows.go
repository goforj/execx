//go:build windows

package execx

import "syscall"

// Setpgid is a no-op on Windows; on Unix it places the child in a new process group.
// @group OS Controls
//
// Example: setpgid
//
//	out, _ := execx.Command("printf", "ok").Setpgid(true).Output()
//	fmt.Print(out)
//	// ok
func (c *Cmd) Setpgid(_ bool) *Cmd {
	return c
}

// Setsid is a no-op on Windows; on Unix it starts a new session.
// @group OS Controls
//
// Example: setsid
//
//	out, _ := execx.Command("printf", "ok").Setsid(true).Output()
//	fmt.Print(out)
//	// ok
func (c *Cmd) Setsid(_ bool) *Cmd {
	return c
}

// Pdeathsig is a no-op on Windows; on Linux it signals the child when the parent exits.
// @group OS Controls
//
// Example: pdeathsig
//
//	out, _ := execx.Command("printf", "ok").Pdeathsig(syscall.SIGTERM).Output()
//	fmt.Print(out)
//	// ok
func (c *Cmd) Pdeathsig(_ syscall.Signal) *Cmd {
	return c
}

// CreationFlags sets Windows process creation flags (for example, create a new process group).
// @group OS Controls
//
// Example: creation flags
//
//	out, _ := execx.Command("printf", "ok").CreationFlags(0x00000200).Output()
//	fmt.Print(out)
//	// ok
func (c *Cmd) CreationFlags(flags uint32) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.CreationFlags = flags
	return c
}

// HideWindow hides console windows and sets CREATE_NO_WINDOW for console apps.
// @group OS Controls
//
// Example: hide window
//
//	out, _ := execx.Command("printf", "ok").HideWindow(true).Output()
//	fmt.Print(out)
//	// ok
func (c *Cmd) HideWindow(on bool) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.HideWindow = on
	if on {
		c.sysProcAttr.CreationFlags |= syscall.CREATE_NO_WINDOW
	}
	return c
}

func (c *Cmd) ensureSysProcAttr() {
	if c.sysProcAttr == nil {
		c.sysProcAttr = &syscall.SysProcAttr{}
	}
}
