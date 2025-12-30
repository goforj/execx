//go:build !windows

package execx

const (
	// CreateNewProcessGroup starts the process in a new process group.
	CreateNewProcessGroup = 0x00000200
	// CreateNewConsole creates a new console for the process.
	CreateNewConsole = 0x00000010
	// CreateNoWindow prevents console windows from being created.
	CreateNoWindow = 0x08000000
)
