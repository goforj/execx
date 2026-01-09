//go:build linux

package execx

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"unsafe"
)

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

func TestPTYIoctlSuccessAndErrorLinux(t *testing.T) {
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		t.Fatalf("open ptmx: %v", err)
	}
	defer master.Close()
	unlock := int32(0)
	if err := ptyIoctl(master.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&unlock))); err != nil {
		t.Fatalf("expected ioctl success, got %v", err)
	}
	if err := ptyIoctl(0, 0, 0); err == nil {
		t.Fatalf("expected ioctl error")
	}
}

func TestOpenPTYWithOpenErrorLinux(t *testing.T) {
	openFile := func(string, int, os.FileMode) (*os.File, error) {
		return nil, errors.New("open failed")
	}
	_, _, err := openPTYWith(openFile, func(uintptr, uintptr, uintptr) error { return nil })
	if err == nil || err.Error() != "open failed" {
		t.Fatalf("expected open error, got %v", err)
	}
}

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

func TestOpenPTYWithSlaveErrorLinux(t *testing.T) {
	openFile := func(name string, flag int, perm os.FileMode) (*os.File, error) {
		if name == "/dev/ptmx" {
			return os.OpenFile(os.DevNull, os.O_RDWR, 0)
		}
		return nil, errors.New("slave open failed")
	}
	ioctl := func(fd uintptr, req uintptr, arg uintptr) error {
		if req == syscall.TIOCGPTN {
			*(*uint32)(unsafe.Pointer(arg)) = 1234
		}
		return nil
	}
	_, _, err := openPTYWith(openFile, ioctl)
	if err == nil || err.Error() != "slave open failed" {
		t.Fatalf("expected slave open error, got %v", err)
	}
}

func TestOpenPTYWithSuccessLinux(t *testing.T) {
	openFile := func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return os.OpenFile(os.DevNull, os.O_RDWR, 0)
	}
	ioctl := func(fd uintptr, req uintptr, arg uintptr) error {
		if req == syscall.TIOCGPTN {
			*(*uint32)(unsafe.Pointer(arg)) = 0
		}
		return nil
	}
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
