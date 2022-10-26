package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"syscall"
)

const expectedOutput = "Gidday world"

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("must pass 2 arguments")
	}

	fd1, err := syscall.Open(os.Args[1], os.O_RDWR|os.O_CREATE|os.O_TRUNC, syscall.S_IRUSR|syscall.S_IWUSR)
	if err != nil {
		return err
	}
	defer syscall.Close(fd1)

	fd2, err := syscall.Dup(fd1)
	if err != nil {
		return err
	}
	defer syscall.Close(fd2)

	fd3, err := syscall.Open(os.Args[1], os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(fd3)

	// Hello,
	if _, err := syscall.Write(fd1, []byte("Hello,")); err != nil {
		return err
	}

	// Hello, world
	if _, err := syscall.Write(fd2, []byte(" world")); err != nil {
		return err
	}

	// Set offest to 0th byte.
	if _, err := syscall.Seek(fd2, 0, os.SEEK_SET); err != nil {
		return err
	}

	// HELLO, world
	if _, err := syscall.Write(fd1, []byte("HELLO")); err != nil {
		return err
	}

	// Gidday world
	if _, err := syscall.Write(fd3, []byte("Gidday")); err != nil {
		return err
	}

	fd4, err := syscall.Open(os.Args[1], os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(fd4)

	output := new(bytes.Buffer)
	buf := make([]byte, 1024)

	for {
		n, err := syscall.Read(fd4, buf);
		if err != nil {
			return err
		}

		if !(n > 0) {
			break
		}

		output.Write(buf[:n])
	}

	if output.String() != expectedOutput {
		return fmt.Errorf("expected %q, got %q", expectedOutput, output.String())
	}

	return nil
}
