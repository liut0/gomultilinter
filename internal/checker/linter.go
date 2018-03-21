package checker

import (
	"fmt"
	"go/token"

	"github.com/liut0/gomultilinter/api"
)

const (
	linterPanciMsg      = "linter did panic: %v"
	linterPanciCategory = "linter-panic"

	linterErrorMsg      = "linter returned error: %v"
	linterErrorCategory = "linter-error"
)

func (c *Checker) lintPkg(pkg *api.Package) {
	for linterName, l := range c.pkgLinter {
		r := c.issueReporter.entry(linterName)
		c.runLinter(r, func() error {
			return l.LintPackage(c.ctx, pkg, r)
		})
	}
}

func (c *Checker) lintFile(file *api.File) {
	for linterName, l := range c.fileLinter {
		r := c.issueReporter.entry(linterName)
		c.runLinter(r, func() error {
			return l.LintFile(c.ctx, file, r)
		})
	}
}

func (c *Checker) runLinter(reporter api.IssueReporter, f func() error) {
	defer func() {
		if err := recover(); err != nil {
			reporter.Report(&api.Issue{
				Message:  fmt.Sprintf(linterPanciMsg, err),
				Position: token.Position{},
				Category: linterPanciCategory,
				Severity: api.SeverityError,
			})
		}

	}()

	if err := f(); err != nil {
		reporter.Report(&api.Issue{
			Message:  fmt.Sprintf(linterErrorMsg, err),
			Position: token.Position{},
			Category: linterErrorCategory,
			Severity: api.SeverityError,
		})
	}
}
