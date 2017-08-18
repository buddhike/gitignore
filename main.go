package main

import (
	"bufio"
	"os"
)

type GitIgnore struct {
	rules []*Rule
}

// Load parses the given .gitignore file and returns a
// list of rules that can be used to evaluate the patterns
// against a provided input.
func Load(path string) (*GitIgnore, error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fs.Close()

	result := []*Rule{}
	scanner := bufio.NewScanner(fs)
	for scanner.Scan() {
		result = append(result, parse(scanner.Text()))
	}

	return &GitIgnore{rules: result}, nil
}

// Match runs the rules in the parsed GitIgnore
// against the specified path.
func (g *GitIgnore) Match(path string) bool {
	shouldIgnore := false
	var rule *Rule
	for _, r := range g.rules {
		rule = r
		m, _ := r.Matcher(newInput(path))
		if m {
			shouldIgnore = !r.IsNegate
		}
	}

	if !shouldIgnore {
		return false
	}

	if rule.IsDir {
		fi, err := os.Stat(path)
		if err != nil {
			return false
		}
		return fi.IsDir()
	}

	return true
}
