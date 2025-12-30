//go:build windows

package execx

import "syscall"

// Setpgid is a no-op on Windows.
// @group OS Controls
//
// Example: setpgid
//
//	fmt.Println(execx.Command("go", "env", "GOOS").Setpgid(true) != nil)
//	// #bool true
func (c *Cmd) Setpgid(_ bool) *Cmd {
	return c
}

// Setsid is a no-op on Windows.
// @group OS Controls
//
// Example: setsid
//
//	fmt.Println(execx.Command("go", "env", "GOOS").Setsid(true) != nil)
//	// #bool true
func (c *Cmd) Setsid(_ bool) *Cmd {
	return c
}

// Pdeathsig is a no-op on Windows.
// @group OS Controls
//
// Example: pdeathsig
//
//	fmt.Println(execx.Command("go", "env", "GOOS").Pdeathsig(0) != nil)
//	// #bool true
func (c *Cmd) Pdeathsig(_ syscall.Signal) *Cmd {
	return c
}

// CreationFlags sets Windows creation flags.
// @group OS Controls
//
// Example: creation flags
//
//	fmt.Println(execx.Command("go", "env", "GOOS").CreationFlags(0) != nil)
//	// #bool true
func (c *Cmd) CreationFlags(flags uint32) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.CreationFlags = flags
	return c
}

// HideWindow controls window visibility and sets CREATE_NO_WINDOW for console apps.
// @group OS Controls
//
// Example: hide window
//
//	fmt.Println(execx.Command("go", "env", "GOOS").HideWindow(true) != nil)
//	// #bool true
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
