package pid

import (
	"syscall"
)

func lock(fd uintptr) error {
	return syscall.Flock(int(fd), syscall.LOCK_EX|syscall.LOCK_NB)
}

func unlock(fd uintptr) error {
	return syscall.Flock(int(fd), syscall.LOCK_UN)
}
