package manager

import (
	"golang.org/x/sys/unix"
)

func SetOSParams() {
	// Prevent coredump to keep process memory confidential
	nocore := &unix.Rlimit{
		Cur: 0,
		Max: 0,
	}

	unix.Setrlimit(unix.RLIMIT_CORE, nocore)
}
