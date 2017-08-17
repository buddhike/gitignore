package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdvance(t *testing.T) {
	i := newInput("a")
	i.advance()
	assert.Equal(t, 1, i.position)

	i.advance()
	assert.Equal(t, 1, i.position)

	i = newInput("")
	i.advance()
	assert.Equal(t, 0, i.position)

	i.advance()
	assert.Equal(t, 0, i.position)
}

func TestCurrent(t *testing.T) {
	i := newInput("ab")

	c, eof := i.current()
	assert.Equal(t, 'a', c)
	assert.False(t, eof)

	i.advance()
	c, eof = i.current()
	assert.Equal(t, 'b', c)
	assert.False(t, eof)

	i.advance()
	c, eof = i.current()
	assert.True(t, eof)
}
