//go:build windows

package execx

import "syscall"

const (
	// CreateNewProcessGroup starts the process in a new process group.
	CreateNewProcessGroup = syscall.CREATE_NEW_PROCESS_GROUP
	// CreateNewConsole creates a new console for the process.
	CreateNewConsole = syscall.CREATE_NEW_CONSOLE
	// CreateNoWindow prevents console windows from being created.
	CreateNoWindow = syscall.CREATE_NO_WINDOW
)
