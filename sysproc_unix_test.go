//go:build unix

package execx

import "testing"

func TestSysProcAttrUnixFlags(t *testing.T) {
	cmd := Command("echo")
	cmd.Setpgid(true).Setsid(true)
	if cmd.sysProcAttr == nil {
		t.Fatalf("expected sys proc attr to be set")
	}
	if !cmd.sysProcAttr.Setpgid {
		t.Fatalf("expected Setpgid to be true")
	}
	if !cmd.sysProcAttr.Setsid {
		t.Fatalf("expected Setsid to be true")
	}

	execCmd := cmd.execCmd()
	if execCmd.SysProcAttr == nil {
		t.Fatalf("expected sys proc attr on exec cmd")
	}
}
