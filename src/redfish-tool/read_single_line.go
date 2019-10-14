package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func readSingleLine(f string) (string, error) {
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
		}
		return line, fmt.Errorf("ERROR: Empty password read from file %s", f)
	}

	return line, nil
}
