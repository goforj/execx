//go:build darwin

package execx

import (
	"bytes"
	"os"
	"syscall"
	"unsafe"
)

// ptyCheck reports Darwin support before command setup allocates any descriptors.
func ptyCheck() error {
	return nil
}

// openPTY delegates system calls so error paths remain deterministic in tests.
func openPTY() (*os.File, *os.File, error) {
	return openPTYWith(os.OpenFile, ptyIoctl)
}

// openPTYWith grants and unlocks a Darwin PTY before opening its discovered slave device.
func openPTYWith(openFile func(string, int, os.FileMode) (*os.File, error), ioctl func(uintptr, uintptr, uintptr) error) (*os.File, *os.File, error) {
	master, err := openFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	if err := ioctl(master.Fd(), syscall.TIOCPTYGRANT, 0); err != nil {
		_ = master.Close()
		return nil, nil, err
	}
	if err := ioctl(master.Fd(), syscall.TIOCPTYUNLK, 0); err != nil {
		_ = master.Close()
		return nil, nil, err
	}
	var nameBuf [128]byte
	if err := ioctl(master.Fd(), syscall.TIOCPTYGNAME, uintptr(unsafe.Pointer(&nameBuf[0]))); err != nil {
		_ = master.Close()
		return nil, nil, err
	}
	name := string(bytes.TrimRight(nameBuf[:], "\x00"))
	slave, err := openFile(name, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		_ = master.Close()
		return nil, nil, err
	}
	return master, slave, nil
}

// ptyIoctl converts the raw syscall errno into an ordinary Go error.
func ptyIoctl(fd uintptr, req uintptr, arg uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, req, arg)
	if errno != 0 {
		return errno
	}
	return nil
}
