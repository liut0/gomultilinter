// Package issue provides models/helpers to the internal representation
// of issues
package issue

import (
	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/internal/files"
)

// LinterIssue is an Issue bound to the linter
// which created it and some additional info
type LinterIssue struct {
	*api.Issue
	Linter string
	Path   Path
}

// Path wraps rel/abs paths
type Path struct {
	Abs string
	Rel string
}

func (p Path) String() string {
	return p.Rel
}

// Line is a shorthand for Position.Line
func (t *LinterIssue) Line() int {
	return t.Position.Line
}

// Col is a shorthand for Position.Column
func (t *LinterIssue) Col() int {
	return t.Position.Column
}

// ToLinterIssue converts an api.Issue to a LinterIssue
func ToLinterIssue(issue *api.Issue, linter string) *LinterIssue {
	return &LinterIssue{
		Issue:  issue,
		Linter: linter,
		Path: Path{
			Abs: files.AbsPath(issue.Position.Filename),
			Rel: files.RelPath(issue.Position.Filename),
		},
	}
}
