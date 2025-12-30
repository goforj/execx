//go:build !windows

package execx

// CreationFlags is a no-op on non-Windows platforms; on Windows it sets process creation flags.
// @group OS Controls
func (c *Cmd) CreationFlags(_ uint32) *Cmd {
	return c
}

// HideWindow is a no-op on non-Windows platforms; on Windows it hides console windows.
// @group OS Controls
func (c *Cmd) HideWindow(_ bool) *Cmd {
	return c
}
