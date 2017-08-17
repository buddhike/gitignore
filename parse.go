package main

const CHAR_SEP = '/'
const CHAR_WILDCARD = '*'
const CHAR_CHOICE_START = '['
const CHAR_CHOICE_END = ']'

type Matcher func(Input) (bool, Input)

// Rule represents a parsed output of a single
// line in .gitignore file.
type Rule struct {
	Matcher Matcher
	IsDir   bool
}

// Helper function to match a given pattern in input.
func seq(pattern string, input Input) (bool, Input) {
	rest := input

	for _, next := range pattern {
		c, eof := rest.current()
		if eof || next != c {
			return false, input
		}
		rest.advance()
	}

	return true, rest
}

// chain takes two Matchers - first and second, then returns a new Matcher
// that returns true when both first and second return true.
func chain(first Matcher, second Matcher) Matcher {
	if first == nil {
		panic("argument nil: first")
	}

	if second == nil {
		panic("argument nil: second")
	}

	return func(i Input) (bool, Input) {
		copy := i
		ok, rest := first(i)
		if !ok {
			return ok, rest
		}

		ok, rest = second(rest)
		if ok {
			return ok, rest
		}
		return false, copy
	}
}

// tryExactMatcher creates a matcher to match each character
// in the pattern until the next marker character.
func tryExactMatcher(pattern Input) (Matcher, Input) {
	p := []rune{}

	for true {
		c, eof := pattern.current()
		if eof {
			break
		}

		p = append(p, c)
		c, eof = pattern.advance()

		if c == CHAR_WILDCARD || c == CHAR_SEP || c == CHAR_CHOICE_START {
			break
		}
	}

	return func(i Input) (bool, Input) {
		return seq(string(p), i)
	}, pattern
}

// tryWildcardMatcher creates a matcher to match
// any character followed by the rest of the pattern string.
// Matching terminates when it encounters the next slash.
func tryWildcardMatcher(pattern Input) (Matcher, Input) {
	ok, rest := seq("*", pattern)
	if !ok {
		return nil, rest
	}

	suffix, rest := createMatcher(rest)

	return func(i Input) (bool, Input) {
		copy := i

		for true {
			ok, rest := suffix(i)
			if ok {
				return ok, rest
			}

			c, eof := i.current()
			if eof || c == CHAR_SEP {
				break
			}

			i.advance()
		}

		return false, copy
	}, rest
}

// tryAnySegmentMatcher creates a matcher to match a slash
// followed by any number of characters followed by an optional
// slash.
// If input only matches the first slash, matcher will return
// true but will consume only first character.
func tryAnySegmentMatcher(i Input) (Matcher, Input) {
	ok, rest := seq("/*/", i)

	if !ok {
		return nil, i
	}

	return func(i Input) (bool, Input) {
		c, _ := i.current()
		if c != CHAR_SEP {
			return false, i
		}

		c, eof := i.advance()
		j := i

		for !eof {
			c, eof = j.current()
			j.advance()

			if c == CHAR_SEP {
				return true, j
			}
		}

		return true, i
	}, rest
}

// tryManySegmentsMatcher creates a matcher to match many slash
// separated segments by rest of the pattern.
// In contrast to tryAnySegmentMatcher, this matcher consumes the slashes.
func tryManySegmentsMatcher(i Input) (Matcher, Input) {
	ok, rest := seq("/**/", i)
	if !ok {
		return nil, i
	}

	suffix, i := createMatcher(rest)

	return func(i Input) (bool, Input) {
		copy := i

		c, eof := i.current()
		if c != CHAR_SEP {
			return false, i
		}
		i.advance()

		for true {
			ok, rest := suffix(i)
			if ok {
				return ok, rest
			}

			c, eof = i.current()
			if eof {
				break
			}
			i.advance()
		}

		return false, copy
	}, i
}

// tryChoiceMatcher creates a matcher to match any character
// in the specified set.
func tryChoiceMatcher(i Input) (Matcher, Input) {
	copy := i
	c, _ := i.current()
	if c != CHAR_CHOICE_START {
		return nil, copy
	}

	choices := make(map[rune]bool)

	for true {
		i.advance()
		c, eof := i.current()
		if eof {
			return nil, copy
		}
		if c == CHAR_CHOICE_END {
			i.advance()
			break
		}
		choices[c] = true
	}

	return func(i Input) (bool, Input) {
		c, _ := i.current()
		if choices[c] {
			i.advance()
			return true, i
		}
		return false, i
	}, i
}

// createMatcher converts an input containing a pattern
// string to a matcher function that can be used to match the
// corresponding pattern.
func createMatcher(i Input) (Matcher, Input) {
	// default matcher returns true without
	// consuming any input.
	p := func(i Input) (bool, Input) {
		return true, i
	}

	for true {
		c, eof := i.current()
		if eof {
			return p, i
		}

		var matcher Matcher
		var rest Input

		switch c {
		case CHAR_SEP:
			matcher, rest = tryManySegmentsMatcher(i)
			if matcher == nil {
				matcher, rest = tryAnySegmentMatcher(i)
			}
		case CHAR_WILDCARD:
			matcher, rest = tryWildcardMatcher(i)
		case CHAR_CHOICE_START:
			matcher, rest = tryChoiceMatcher(i)
		}

		if matcher == nil {
			matcher, rest = tryExactMatcher(i)
		}

		p = chain(p, matcher)
		i = rest
	}

	return p, i
}

func parse(line string) *Rule {
	i := newInput(line)
	p, _ := createMatcher(i)
	l, _ := i.last()

	return &Rule{
		Matcher: p,
		IsDir:   l == CHAR_SEP,
	}
}
