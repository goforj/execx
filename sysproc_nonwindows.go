//go:build !windows

package execx

// CreationFlags is a no-op on non-Windows platforms.
// @group OS Controls
//
// Example: creation flags
//
//	fmt.Println(execx.Command("go", "env", "GOOS").CreationFlags(0) != nil)
//	// #bool true
func (c *Cmd) CreationFlags(_ uint32) *Cmd {
	return c
}

// HideWindow is a no-op on non-Windows platforms.
// @group OS Controls
//
// Example: hide window
//
//	fmt.Println(execx.Command("go", "env", "GOOS").HideWindow(true) != nil)
//	// #bool true
func (c *Cmd) HideWindow(_ bool) *Cmd {
	return c
}
