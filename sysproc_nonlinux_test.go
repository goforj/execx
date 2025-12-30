//go:build unix && !linux

package execx

import (
	"syscall"
	"testing"
)

func TestPdeathsigNoop(t *testing.T) {
	cmd := Command("echo")
	cmd.Pdeathsig(syscall.SIGTERM)
	if cmd.sysProcAttr != nil {
		t.Fatalf("expected Pdeathsig to be a no-op")
	}
}
