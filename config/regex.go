package config

import (
	"regexp"
)

// Regex is an JSON-Compatible wrapper for a compiled regular expression
type Regex struct {
	*regexp.Regexp
}

// MultiRegex wraps multiple regular expressions
type MultiRegex []*Regex

// MatchesAny returns true if any of the regexps matches the input
func (r MultiRegex) MatchesAny(input string) bool {
	for _, rgx := range r {
		if rgx.MatchString(input) {
			return true
		}
	}
	return false
}

// MatchesAll returns true if all of the regexps matches the input
func (r MultiRegex) MatchesAll(input string) bool {
	for _, rgx := range r {
		if !rgx.MatchString(input) {
			return false
		}
	}
	return true
}

// UnmarshalText compiles the provided regex
func (r *Regex) UnmarshalText(data []byte) error {
	var err error
	r.Regexp, err = regexp.Compile(string(data))
	return err
}
