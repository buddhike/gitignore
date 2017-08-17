package main

import (
	"bufio"
	"os"
)

func loadGitIgnore(path string) ([]string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []string{}, nil
	}

	fd, err := os.Open(path)
	defer fd.Close()

	if err != nil {
		return nil, err
	}

	scan := bufio.NewScanner(fd)
	lines := make([]string, 16)

	for scan.Scan() {
		lines = append(lines, scan.Text())
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return lines, err
}
