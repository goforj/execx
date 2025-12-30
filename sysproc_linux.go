//go:build linux

package execx

import "syscall"

// Setpgid sets the process group ID behavior.
// @group OS Controls
//
// Example: setpgid
//
//	fmt.Println(execx.Command("go", "env", "GOOS").Setpgid(true) != nil)
//	// #bool true
func (c *Cmd) Setpgid(on bool) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.Setpgid = on
	return c
}

// Setsid sets the session ID behavior.
// @group OS Controls
//
// Example: setsid
//
//	fmt.Println(execx.Command("go", "env", "GOOS").Setsid(true) != nil)
//	// #bool true
func (c *Cmd) Setsid(on bool) *Cmd {
	c.ensureSysProcAttr()
	c.sysProcAttr.Setsid = on
	return c
}

// Pdeathsig sets a parent-death signal on Linux.
// @group OS Controls
//
// Example: pdeathsig
//
//	fmt.Println(execx.Command("go", "env", "GOOS").Pdeathsig(0) != nil)
//	// #bool true
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
