package main

import (
	"fmt"
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
	if len(os.Args) < 3 {
		return fmt.Errorf("%s SRC DEST\n", os.Args[0])
	}

	src, err := syscall.Open(os.Args[1], os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer syscall.Close(src)

	dst, err := syscall.Open(os.Args[2], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer syscall.Close(dst)

	buf := make([]byte, 1024)

	for {
		n, err := syscall.Read(src, buf)
		if err != nil {
			return err
		}

		if !(n > 0) {
			break
		}

		_, err = syscall.Write(dst, buf[0:n])
		if err != nil {
			return err
		}
	}

	return nil
}
