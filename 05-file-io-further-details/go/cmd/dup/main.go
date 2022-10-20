package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	fd, err := dup(uintptr(syscall.Stdout))
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf("dup() -> %d\n", fd))
	_, err = syscall.Write(fd, buf.Bytes())
	if err != nil {
		return err
	}

	fd2, err := dup2(uintptr(fd), uintptr(fd)+5)
	if err != nil {
		return err
	}
	defer syscall.Close(fd2)

	buf.Reset()

	buf.WriteString(fmt.Sprintf("dup2() -> %d\n", fd2))
	_, err = syscall.Write(fd2, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// dup is an implementation of the dup() system call using fnctl(). It returns
// a duplicate of the file descriptor oldfd.
func dup(oldfd uintptr) (int, error) {
	_, err := unix.FcntlInt(oldfd, unix.F_GETFD, 0)
	if err != nil {
		return 0, err
	}

	return unix.FcntlInt(oldfd, unix.F_DUPFD, 0)
}

// dup2 is an implementation of the dup2() system call using fnctl(). It
// returns a duplicate of the file descriptor oldfd using newfd as the new file
// descriptor.
func dup2(oldfd uintptr, newfd uintptr) (int, error) {
	_, err := unix.FcntlInt(oldfd, unix.F_GETFD, 0)
	if err != nil {
		return 0, err
	}

	// When both file descriptors match, check that the old file descriptor is
	// valid and return the new file descriptor if so.
	if oldfd == newfd {
		_, err := unix.FcntlInt(oldfd, unix.F_GETFL, 0)
		if err != nil {
			return -1, syscall.EBADF
		}
		return int(newfd), nil
	}

	// Check if the new file descriptor is valid.
	_, err = unix.FcntlInt(newfd, unix.F_GETFD, 0)
	if err != nil {
		if !errors.Is(err, syscall.EBADF) {
			return -1, err
		}
	}

	// Make a best effort to close the new file descriptor.
	syscall.Close(int(newfd))

	return unix.FcntlInt(oldfd, unix.F_DUPFD, int(newfd))
}
