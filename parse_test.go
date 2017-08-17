package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeq(t *testing.T) {
	s, r := seq("", newInput(""))
	assert.True(t, s)

	p := "abc"
	s, r = seq(p, newInput("abc"))
	assert.True(t, s)
	assert.Equal(t, 3, r.position)

	s, r = seq(p, newInput("abcd"))
	assert.True(t, s)
	assert.Equal(t, 3, r.position)

	s, r = seq(p, newInput("ab"))
	assert.False(t, s)
	assert.Equal(t, 0, r.position)

	s, r = seq(p, newInput(""))
	assert.False(t, s)
	assert.Equal(t, 0, r.position)
}

type expectation struct {
	Input    string
	Result   bool
	Position int
}

type testPattern struct {
	Pattern      string
	Expectations []expectation
}

func TestMatch(t *testing.T) {
	tests := []testPattern{
		{
			Pattern: "abc",
			Expectations: []expectation{
				{"abc", true, 3},
				{"abcabc", true, 3},
				{"abc/", true, 3},
				{"", false, 0},
				{"a", false, 0},
				{"abx", false, 0},
				{"abxy", false, 0},
			},
		},
		{
			Pattern: "a*c",
			Expectations: []expectation{
				{"abc", true, 3},
				{"abc/", true, 3},
				{"a", false, 0},
				{"abx", false, 0},
				{"abxy", false, 0},
				{"ab/c", false, 0},
			},
		},
		{
			Pattern: "a/**/c",
			Expectations: []expectation{
				{"a/c", true, 3},
				{"a/b/c", true, 5},
				{"a/b/d/c", true, 7},
				{"a/d", false, 0},
				{"a/b/d", false, 0},
				{"abc", false, 0},
			},
		},
		{
			Pattern: "a/*/c",
			Expectations: []expectation{
				{"a/c", true, 3},
				{"a/b/c", true, 5},
				{"a/b/c/d", true, 5},
				{"a/d", false, 0},
				{"a/b/d", false, 0},
			},
		},
		{
			Pattern: "a/*/",
			Expectations: []expectation{
				{"a/c", true, 2},
				{"a/c/d", true, 4},
			},
		},
		{
			Pattern: "a/*b*/c",
			Expectations: []expectation{
				{"a/b/c", true, 5},
				{"a/aba/c", true, 7},
				{"a/bbb/c", true, 7},
				{"a/aba/c/e/f/g", true, 7},
				{"a/c", false, 0},
				{"a/aaa/c", false, 0},
			},
		},
		{
			Pattern: "a/[abc]d/d",
			Expectations: []expectation{
				{"a/ad/d", true, 6},
				{"a/bd/d", true, 6},
				{"a/dd/d", false, 0},
			},
		},
	}

	for _, test := range tests {
		p := parse(test.Pattern)
		for _, expectation := range test.Expectations {
			ok, rest := p.Matcher(newInput(expectation.Input))
			assert.Equal(t, expectation.Result, ok, test.Pattern, expectation.Input)
			assert.Equal(t, expectation.Position, rest.position, test.Pattern, expectation.Input)
		}
	}
}
