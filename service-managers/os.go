package manager

import "syscall"

func SetOSParams() {
	// Prevent coredump to keep process memory confidential
	nocore := &syscall.Rlimit{
		Cur: 0,
		Max: 0,
	}

	syscall.Setrlimit(syscall.RLIMIT_CORE, nocore)
}
