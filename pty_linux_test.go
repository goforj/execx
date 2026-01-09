//go:build linux

package execx

import "testing"

func TestPTYLinuxOpen(t *testing.T) {
	if err := ptyCheck(); err != nil {
		t.Fatalf("unexpected pty check error: %v", err)
	}
	master, slave, err := openPTY()
	if err != nil {
		t.Fatalf("expected openPTY to succeed, got %v", err)
	}
	_ = master.Close()
	_ = slave.Close()
}
