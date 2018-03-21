package config

import (
	"strings"

	"github.com/liut0/gomultilinter/api"
)

// Severity is an JSON-compatible severity
type Severity struct {
	api.Severity
}

// UnmarshalText parses the severity out of the provided text
func (s *Severity) UnmarshalText(data []byte) error {
	var err error
	s.Severity, err = api.ParseSeverity(strings.Title(string(data)))
	return err
}
