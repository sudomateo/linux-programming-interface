package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	appendMode := flag.Bool("a", false, "append to given FILEs, do not overwrite")
	flag.Parse()

	fileMode := os.O_CREATE | os.O_WRONLY
	if *appendMode {
		fileMode = fileMode | os.O_APPEND
	} else {
		fileMode = fileMode | os.O_TRUNC
	}

	files := make([]int, 0)

	// Cleanup opened files upon exit.
	defer func() {
		for _, file := range files {
			_ = syscall.Close(file)
		}
	}()

	fileArgs := flag.Args()
	if len(fileArgs) > 0 {
		for _, file := range fileArgs {
			fd, err := syscall.Open(file, fileMode, 0644)
			if err != nil {
				return fmt.Errorf("opening file %q: %v", file, err)
			}

			files = append(files, fd)
		}
	}

	buf := make([]byte, 1024)

	for {
		n, err := syscall.Read(syscall.Stdin, buf)
		if err != nil {
			return fmt.Errorf("reading from stdin: %v", err)
		}

		// There's nothing left to read on stdout.
		if !(n > 0) {
			break
		}

		_, err = syscall.Write(syscall.Stdout, buf[0:n])
		if err != nil {
			return fmt.Errorf("writing to stdout: %v", err)
		}

		for _, fd := range files {
			_, err = syscall.Write(fd, buf[0:n])
			if err != nil {
				return fmt.Errorf("writing to file descriptor %d: %v", fd, err)
			}
		}
	}

	return nil
}
