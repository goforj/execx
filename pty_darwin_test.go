//go:build darwin

package execx

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"unsafe"
)

func TestPTYDarwinOpen(t *testing.T) {
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

func TestPTYIoctlSuccessAndError(t *testing.T) {
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		t.Fatalf("open ptmx: %v", err)
	}
	defer master.Close()
	if err := ptyIoctl(master.Fd(), syscall.TIOCPTYGRANT, 0); err != nil {
		t.Fatalf("expected ioctl success, got %v", err)
	}
	if err := ptyIoctl(0, 0, 0); err == nil {
		t.Fatalf("expected ioctl error")
	}
}

func TestOpenPTYWithOpenError(t *testing.T) {
	openFile := func(string, int, os.FileMode) (*os.File, error) {
		return nil, errors.New("open failed")
	}
	_, _, err := openPTYWith(openFile, func(uintptr, uintptr, uintptr) error { return nil })
	if err == nil || err.Error() != "open failed" {
		t.Fatalf("expected open error, got %v", err)
	}
}

func TestOpenPTYWithGrantError(t *testing.T) {
	openFile := func(string, int, os.FileMode) (*os.File, error) {
		return os.OpenFile(os.DevNull, os.O_RDWR, 0)
	}
	_, _, err := openPTYWith(openFile, func(fd uintptr, req uintptr, arg uintptr) error {
		if req == syscall.TIOCPTYGRANT {
			return errors.New("grant failed")
		}
		return nil
	})
	if err == nil || err.Error() != "grant failed" {
		t.Fatalf("expected grant error, got %v", err)
	}
}

func TestOpenPTYWithUnlockError(t *testing.T) {
	openFile := func(string, int, os.FileMode) (*os.File, error) {
		return os.OpenFile(os.DevNull, os.O_RDWR, 0)
	}
	ioctl := func(fd uintptr, req uintptr, arg uintptr) error {
		if req == syscall.TIOCPTYUNLK {
			return errors.New("unlock failed")
		}
		return nil
	}
	_, _, err := openPTYWith(openFile, ioctl)
	if err == nil || err.Error() != "unlock failed" {
		t.Fatalf("expected unlock error, got %v", err)
	}
}

func TestOpenPTYWithNameError(t *testing.T) {
	openFile := func(string, int, os.FileMode) (*os.File, error) {
		return os.OpenFile(os.DevNull, os.O_RDWR, 0)
	}
	ioctl := func(fd uintptr, req uintptr, arg uintptr) error {
		if req == syscall.TIOCPTYGNAME {
			return errors.New("name failed")
		}
		return nil
	}
	_, _, err := openPTYWith(openFile, ioctl)
	if err == nil || err.Error() != "name failed" {
		t.Fatalf("expected name error, got %v", err)
	}
}

func TestOpenPTYWithSlaveError(t *testing.T) {
	openFile := func(name string, flag int, perm os.FileMode) (*os.File, error) {
		if name == "/dev/ptmx" {
			return os.OpenFile(os.DevNull, os.O_RDWR, 0)
		}
		return nil, errors.New("slave open failed")
	}
	ioctl := func(fd uintptr, req uintptr, arg uintptr) error {
		if req == syscall.TIOCPTYGNAME {
			buf := (*[128]byte)(unsafe.Pointer(arg))
			copy(buf[:], []byte("/dev/doesnotexist"))
		}
		return nil
	}
	_, _, err := openPTYWith(openFile, ioctl)
	if err == nil || err.Error() != "slave open failed" {
		t.Fatalf("expected slave open error, got %v", err)
	}
}

func TestOpenPTYWithSuccess(t *testing.T) {
	openFile := func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return os.OpenFile(os.DevNull, os.O_RDWR, 0)
	}
	ioctl := func(fd uintptr, req uintptr, arg uintptr) error {
		if req == syscall.TIOCPTYGNAME {
			buf := (*[128]byte)(unsafe.Pointer(arg))
			copy(buf[:], []byte(os.DevNull))
		}
		return nil
	}
	master, slave, err := openPTYWith(openFile, ioctl)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	_ = master.Close()
	_ = slave.Close()
}
