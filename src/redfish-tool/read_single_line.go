package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func ReadSingleLine(f string) (string, error) {
	var line string
	fd, err := os.Open(f)
	if err != nil {
		return line, err
	}

	scanner := bufio.NewScanner(fd)
	scanner.Scan()
	line = scanner.Text()
	fd.Close()
	if line == "" {
		return line, errors.New(fmt.Sprintf("ERROR: Empty password read from file %s", f))
	}

	return line, nil
}
