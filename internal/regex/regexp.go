// Package regex utils to process regular expressions
package regex

import (
	"regexp"
)

// NiceRegexp is a wrapper around *regexp.Regexp
// which provides additional helper funcs
type NiceRegexp struct {
	*regexp.Regexp
}

// MustCompile is a shortcut for regexp.MustCompile wrapped in NiceRegexp
func MustCompile(s string) NiceRegexp {
	return NiceRegexp{Regexp: regexp.MustCompile(s)}
}

// FindNamedStringSubmatch uses regexp.FindStringSubmatch and maps non empty
// named matched groups to the result map
// A return value of nil indicates no match.
func (r NiceRegexp) FindNamedStringSubmatch(s string) map[string]string {
	match := r.FindStringSubmatch(s)
	if match == nil {
		return nil
	}

	result := map[string]string{}
	for i, name := range r.SubexpNames() {
		m := match[i]
		if len(name) > 0 && len(m) > 0 {
			result[name] = m
		}
	}

	return result
}
