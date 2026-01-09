//go:build linux

package execx

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func ptyCheck() error {
	return nil
}

func openPTY() (*os.File, *os.File, error) {
	return openPTYWith(os.OpenFile, ptyIoctl)
}

func openPTYWith(openFile func(string, int, os.FileMode) (*os.File, error), ioctl func(uintptr, uintptr, uintptr) error) (*os.File, *os.File, error) {
	master, err := openFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	fd := master.Fd()
	unlock := int32(0)
	if err := ioctl(fd, syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&unlock))); err != nil {
		_ = master.Close()
		return nil, nil, err
	}
	var ptyNum uint32
	if err := ioctl(fd, syscall.TIOCGPTN, uintptr(unsafe.Pointer(&ptyNum))); err != nil {
		_ = master.Close()
		return nil, nil, err
	}
	name := fmt.Sprintf("/dev/pts/%d", ptyNum)
	slave, err := openFile(name, os.O_RDWR|syscall.O_NOCTTY, 0)
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
