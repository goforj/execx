//go:build !windows

package execx

// CreationFlags is a no-op on non-Windows platforms; on Windows it sets process creation flags.
// @group OS Controls
//
// Example: creation flags
//
//	out, _ := execx.Command("printf", "ok").CreationFlags(0x00000200).Output()
//	fmt.Print(out)
//	// ok
func (c *Cmd) CreationFlags(_ uint32) *Cmd {
	return c
}

// HideWindow is a no-op on non-Windows platforms; on Windows it hides console windows.
// @group OS Controls
//
// Example: hide window
//
//	out, _ := execx.Command("printf", "ok").HideWindow(true).Output()
//	fmt.Print(out)
//	// ok
func (c *Cmd) HideWindow(_ bool) *Cmd {
	return c
}
