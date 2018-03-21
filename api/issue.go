package api

import "go/token"

// IssueReporter collects resultig issues from the linters
// fields is a list of key-value-pairs to include in the log msg
type IssueReporter interface {
	Debug(msg string, fields ...interface{})
	Report(issue *Issue)
}

// Issue represents a linter-issue
type Issue struct {
	Position token.Position
	Severity Severity
	Category string
	Message  string
}
