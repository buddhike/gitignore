package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestFile(entries []string) error {
	os.Mkdir(".tmp", 0755)
	path := ".tmp/test.gitignore"
	os.Remove(path)
	fs, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer fs.Close()

	for _, v := range entries {
		fs.WriteString(fmt.Sprintf("%s\n", v))
	}
	fs.Sync()
	return nil
}

func TestLoad(t *testing.T) {
	err := createTestFile([]string{"abc", "d*f", "a/**/b"})
	if err != nil {
		t.Fatal(err)
	}

	r, le := Load(".tmp/test.gitignore")
	if le != nil {
		t.Fatal(le)
	}

	assert.True(t, r.Match("abc"))
}

func TestNegate(t *testing.T) {
	err := createTestFile([]string{"abc", "!a"})
	if err != nil {
		t.Fatal(err)
	}

	r, le := Load(".tmp/test.gitignore")
	if le != nil {
		t.Fatal(le)
	}

	ok := r.Match("abc")
	assert.False(t, ok)
}
