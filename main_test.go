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
	err := createTestFile([]string{"abc/de*", "!abc/def"})
	if err != nil {
		t.Fatal(err)
	}

	r, le := Load(".tmp/test.gitignore")
	if le != nil {
		t.Fatal(le)
	}

	assert.True(t, r.Match("abc/deh"))
	assert.False(t, r.Match("abc/def"))
}

func TestRules(t *testing.T) {
	i := parse("abc")
	assert.True(t, i.Matcher("abc"))
	assert.True(t, i.Matcher("abc/"))
	assert.False(t, i.Matcher("abc/ "))
	assert.False(t, i.Matcher("abcabc"))
	assert.False(t, i.Matcher(""))
	assert.False(t, i.Matcher("a"))
}

func TestWildcardInBetween(t *testing.T) {
	i := parse("a*c")
	assert.True(t, i.Matcher("abc"))
	assert.True(t, i.Matcher("abc/"))
	assert.True(t, i.Matcher("aaac"))
	assert.False(t, i.Matcher("abcd"))
	assert.False(t, i.Matcher("xabc"))
}

func TestMultipleDirectoryMatcher(t *testing.T) {
	i := parse("a/**/c")
	assert.True(t, i.Matcher("a/c"))
	assert.True(t, i.Matcher("a/b/c"))
	assert.True(t, i.Matcher("a/b/c/"))
	assert.True(t, i.Matcher("a/b/d/c"))
	assert.False(t, i.Matcher("a/d"))
	assert.False(t, i.Matcher("ab/c"))
	assert.False(t, i.Matcher("a/b/c/d"))
}

func TestAnyDirectoryMatcher(t *testing.T) {
	i := parse("a/*/c")
	assert.True(t, i.Matcher("a/c"))
	assert.True(t, i.Matcher("a/b/c"))
	assert.False(t, i.Matcher("a/b/d/c"))
	assert.False(t, i.Matcher("ab/c"))
}

func TestAllTrailingDirectories(t *testing.T) {
	i := parse("a/*/")
	assert.False(t, i.Matcher("a/c"))
	assert.False(t, i.Matcher("a/c/d"))
}

func TestNegateRule(t *testing.T) {
	i := parse("!abc")
	assert.True(t, i.Matcher("abc"))
}

func TestTrailingSlash(t *testing.T) {
	i := parse("abc/")
	assert.True(t, i.Matcher("abc/"))
}
