package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
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
		return fmt.Errorf("%s FILE NUM_BYTES [x]", os.Args[0])
	}

	numBytes, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}

	var skipAppend bool
	fileMode := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if len(os.Args) > 3 && os.Args[3] == "x" {
		skipAppend = true
		fileMode = fileMode ^ os.O_APPEND
	}

	fd, err := syscall.Open(os.Args[1], fileMode, 0644)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	for i := 0; i < numBytes; i++ {
		if skipAppend {
			_, err = syscall.Seek(fd, 0, io.SeekEnd)
			if err != nil {
				return err
			}
		}

		_, err = syscall.Write(fd, []byte("X"))
		if err != nil {
			return err
		}
	}

	return nil
}
