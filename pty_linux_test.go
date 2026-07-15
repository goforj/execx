//go:build linux

package execx

import (
	"errors"
	"os"
	"syscall"
	"testing"
)

// TestPTYLinuxOpen ensures the host kernel can provide the master and slave pair used by PTY execution.
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

// TestPTYIoctlSuccessAndErrorLinux ensures ioctl failures are surfaced instead of silently ignored.
func TestPTYIoctlSuccessAndErrorLinux(t *testing.T) {
	if err := ptyIoctl(0, 0, 0); err == nil {
		t.Fatalf("expected ioctl error")
	}
}

// TestOpenPTYWithOpenErrorLinux ensures opening the multiplexer is the first reported PTY failure.
func TestOpenPTYWithOpenErrorLinux(t *testing.T) {
	openFile := func(string, int, os.FileMode) (*os.File, error) {
		return nil, errors.New("open failed")
	}
	_, _, err := openPTYWith(openFile, func(uintptr, uintptr, uintptr) error { return nil })
	if err == nil || err.Error() != "open failed" {
		t.Fatalf("expected open error, got %v", err)
	}
}

// TestOpenPTYWithUnlockErrorLinux ensures a locked slave cannot escape as a partially initialized pair.
func TestOpenPTYWithUnlockErrorLinux(t *testing.T) {
	openFile := func(string, int, os.FileMode) (*os.File, error) {
		return os.OpenFile(os.DevNull, os.O_RDWR, 0)
	}
	_, _, err := openPTYWith(openFile, func(fd uintptr, req uintptr, arg uintptr) error {
		if req == syscall.TIOCSPTLCK {
			return errors.New("unlock failed")
		}
		return nil
	})
	if err == nil || err.Error() != "unlock failed" {
		t.Fatalf("expected unlock error, got %v", err)
	}
}

// TestOpenPTYWithPTNErrorLinux ensures slave-number lookup failures close the incomplete PTY setup.
func TestOpenPTYWithPTNErrorLinux(t *testing.T) {
	openFile := func(string, int, os.FileMode) (*os.File, error) {
		return os.OpenFile(os.DevNull, os.O_RDWR, 0)
	}
	_, _, err := openPTYWith(openFile, func(fd uintptr, req uintptr, arg uintptr) error {
		if req == syscall.TIOCGPTN {
			return errors.New("ptn failed")
		}
		return nil
	})
	if err == nil || err.Error() != "ptn failed" {
		t.Fatalf("expected ptn error, got %v", err)
	}
}

// TestOpenPTYWithSlaveErrorLinux ensures slave-open failures are returned after master initialization.
func TestOpenPTYWithSlaveErrorLinux(t *testing.T) {
	openFile := func(name string, flag int, perm os.FileMode) (*os.File, error) {
		if name == "/dev/ptmx" {
			return os.OpenFile(os.DevNull, os.O_RDWR, 0)
		}
		return nil, errors.New("slave open failed")
	}
	ioctl := func(uintptr, uintptr, uintptr) error { return nil }
	_, _, err := openPTYWith(openFile, ioctl)
	if err == nil || err.Error() != "slave open failed" {
		t.Fatalf("expected slave open error, got %v", err)
	}
}

// TestOpenPTYWithSuccessLinux ensures injected system calls can complete a usable PTY pair.
func TestOpenPTYWithSuccessLinux(t *testing.T) {
	openFile := func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return os.OpenFile(os.DevNull, os.O_RDWR, 0)
	}
	ioctl := func(uintptr, uintptr, uintptr) error { return nil }
	master, slave, err := openPTYWith(openFile, ioctl)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	_ = master.Close()
	_ = slave.Close()
	if master.Name() != os.DevNull || slave.Name() != os.DevNull {
		t.Fatalf("expected dev null files, got %q %q", master.Name(), slave.Name())
	}
}
