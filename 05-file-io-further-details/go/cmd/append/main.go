package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("%s FILE", os.Args[0])
	}

	fd, err := syscall.Open(os.Args[1], os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	_, err = syscall.Seek(fd, 0, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = syscall.Write(fd, []byte("Hello, World!"))
	if err != nil {
		return err
	}

	return nil
}
