//go:build !windows

package execx

// CreationFlags is a no-op on non-Windows platforms.
func (c *Cmd) CreationFlags(_ uint32) *Cmd {
	return c
}

// HideWindow is a no-op on non-Windows platforms.
func (c *Cmd) HideWindow(_ bool) *Cmd {
	return c
}
