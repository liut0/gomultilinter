package filter

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/internal/checker/issue"
	"github.com/liut0/gomultilinter/internal/regex"
)

const (
	noLinterRgxGrpLinter  = "LINTER"
	noLinterRgxLinterName = `[A-Za-z0-9_\-]+`

	categoryUnnecessaryNoLinterDirective = "unnecessary-nolinter-directive"
	msgUnnecessaryNoLinterDirective      = "unnecessary nolinter directive detected"
)

var (
	noLinterRgx = regex.MustCompile(
		`// ?nolint(: )?` +
			`(?P<` + noLinterRgxGrpLinter + `>(` + noLinterRgxLinterName + `)(, ?` + noLinterRgxLinterName + `)*)?`)
)

// NoLinterDirectiveFilter filters out issues to which a nolinter directive applies
type NoLinterDirectiveFilter struct {
	disabled bool

	fset        *token.FileSet
	parsedFiles map[string]bool

	ranges []*noLinterRange
}

type noLinterRange struct {
	linters      map[string]bool
	filterLinter bool

	start token.Position
	end   token.Position

	path string

	necessary bool
}

// AddFile indexes all nolint directives in this file
func (f *NoLinterDirectiveFilter) AddFile(file *api.File) {
	for node, cmntGrps := range file.CommentMap {
		for _, cmntGrp := range cmntGrps {
			for _, cmnt := range cmntGrp.List {
				f.parseNoLinterRange(cmnt, file, node)
			}
		}
	}
}

func (f *NoLinterDirectiveFilter) parseNoLinterRange(cmnt *ast.Comment, file *api.File, affectedNode ast.Node) {
	match := noLinterRgx.FindNamedStringSubmatch(cmnt.Text)
	if match == nil {
		return
	}

	linters := map[string]bool{}

	linterNames, ok := match[noLinterRgxGrpLinter]
	if ok {
		for _, linterName := range strings.Split(linterNames, ",") {
			linters[strings.TrimSpace(linterName)] = true
		}
	}

	f.ranges = append(f.ranges, &noLinterRange{
		path:         file.Position.Filename,
		start:        file.FSet.Position(affectedNode.Pos()),
		end:          file.FSet.Position(affectedNode.End()),
		linters:      linters,
		filterLinter: len(linters) > 0,
	})
}

// IgnoreIssue returns wether this issue should be ignored (true) or written out (false)
func (f *NoLinterDirectiveFilter) IgnoreIssue(issue *issue.LinterIssue) bool {
	if f.disabled {
		return false
	}

	for _, r := range f.ranges {
		if r.includes(issue) {
			return true
		}
	}
	return false
}

// ReportUnnecessaryDirectives reports unnecessary nolint directives as issues
func (f *NoLinterDirectiveFilter) ReportUnnecessaryDirectives(reporter api.IssueReporter) {
	f.disabled = true
	for _, r := range f.ranges {
		if !r.necessary {
			reporter.Report(&api.Issue{
				Position: r.start,
				Severity: api.SeverityWarning,
				Category: categoryUnnecessaryNoLinterDirective,
				Message:  msgUnnecessaryNoLinterDirective,
			})
		}
	}
	f.disabled = false
}

func (r *noLinterRange) includes(issue *issue.LinterIssue) bool {
	if r.path != issue.Path.Abs || issue.Line() < r.start.Line || issue.Line() > r.end.Line {
		return false
	}

	if !r.filterLinter || r.linters[issue.Linter] {
		r.necessary = true
		return true
	}

	return false
}
