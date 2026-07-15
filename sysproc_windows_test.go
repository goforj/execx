//go:build windows

package execx

import "testing"

// TestHideWindowSetsCreateNoWindow ensures window suppression preserves caller-supplied creation flags.
func TestHideWindowSetsCreateNoWindow(t *testing.T) {
	cmd := Command("go", "env", "GOOS")
	cmd.CreationFlags(0x10).HideWindow(true)
	if cmd.sysProcAttr == nil {
		t.Fatalf("expected sys proc attr")
	}
	if !cmd.sysProcAttr.HideWindow {
		t.Fatalf("expected HideWindow set")
	}
	if cmd.sysProcAttr.CreationFlags&CreateNoWindow == 0 {
		t.Fatalf("expected CREATE_NO_WINDOW flag")
	}
	if cmd.sysProcAttr.CreationFlags&0x10 == 0 {
		t.Fatalf("expected custom creation flags preserved")
	}
}
