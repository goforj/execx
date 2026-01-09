//go:build darwin

package execx

import (
	"bytes"
	"os"
	"syscall"
	"unsafe"
)

func ptyCheck() error {
	return nil
}

func openPTY() (*os.File, *os.File, error) {
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}
	var nameBuf [128]byte
	if err := ptyIoctl(master.Fd(), syscall.TIOCPTYGNAME, uintptr(unsafe.Pointer(&nameBuf[0]))); err != nil {
		_ = master.Close()
		return nil, nil, err
	}
	name := string(bytes.TrimRight(nameBuf[:], "\x00"))
	slave, err := os.OpenFile(name, os.O_RDWR, 0)
	if err != nil {
		_ = master.Close()
		return nil, nil, err
	}
	return master, slave, nil
}

func ptyIoctl(fd uintptr, req uintptr, arg uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, req, arg)
	if errno != 0 {
		return errno
	}
	return nil
}
