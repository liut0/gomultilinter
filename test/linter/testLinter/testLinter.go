package main

import (
	"context"
	"errors"

	"github.com/liut0/gomultilinter/api"
)

const (
	category = "test"
	msg      = "testmsg"
)

// TestLinter does some dummy linting
type TestLinter struct{}

// LinterFactory exported plugin symb
var LinterFactory api.LinterFactory = &TestLinter{}

// NewLinterConfig constructs new config struct
func (l *TestLinter) NewLinterConfig() api.LinterConfig {
	return &TestLinter{}
}

// NewLinter constructs new linter
func (l *TestLinter) NewLinter() (api.Linter, error) {
	return l, nil
}

// Name of the linter
func (*TestLinter) Name() string {
	return "testlinter"
}

// LintFile does some dummy linting based on comments
func (l *TestLinter) LintFile(ctx context.Context, file *api.File, reporter api.IssueReporter) error {
	for _, cmntGrp := range file.ASTFile.Comments {
		for _, cmnt := range cmntGrp.List {
			switch cmnt.Text {
			case "// issue":
				reporter.Report(&api.Issue{
					Severity: api.SeverityWarning,
					Category: category,
					Position: file.FSet.Position(cmnt.Pos()),
					Message:  msg,
				})
			case "// ignored":
				reporter.Report(&api.Issue{
					Severity: api.SeverityWarning,
					Category: "ignored-category",
					Position: file.FSet.Position(cmnt.Pos()),
					Message:  "ignored-msg",
				})
			case "// panic":
				panic(msg)
			case "// error":
				return errors.New(msg)
			}
		}
	}

	return nil
}
