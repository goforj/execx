//go:build unix && !linux

package execx

import (
	"syscall"
	"testing"
)

// TestPdeathsigNoop ensures the Linux-only parent-death signal API remains harmless on other Unix systems.
func TestPdeathsigNoop(t *testing.T) {
	cmd := Command("echo")
	cmd.Pdeathsig(syscall.SIGTERM)
	if cmd.sysProcAttr != nil {
		t.Fatalf("expected Pdeathsig to be a no-op")
	}
}
