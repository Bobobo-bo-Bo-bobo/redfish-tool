package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func ReadSingleLine(f string) (string, error) {
	var line string
	var fd *os.File
	var err error

	if f == "-" {
		fd = os.Stdin
	} else {
		fd, err = os.Open(f)
		if err != nil {
			return line, err
		}
	}

	scanner := bufio.NewScanner(fd)
	scanner.Scan()
	line = scanner.Text()
	fd.Close()
	if line == "" {
		if f == "-" {
			return line, errors.New("ERROR: Empty password read from stdin")
		} else {
			return line, errors.New(fmt.Sprintf("ERROR: Empty password read from file %s", f))
		}
	}

	return line, nil
}
