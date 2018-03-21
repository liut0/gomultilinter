package filter

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/internal/checker/issue"
	"github.com/liut0/gomultilinter/internal/files"
	"github.com/stretchr/testify/assert"
)

const testSrc = `// just a random cmnt
package p

// nolint
func foo() {

}

func bar() int {
	return 1// nolint: foolinter, barlinter
}

// nolint
func bar2() int {
	return 1
}`

type testReporter struct {
	issues []*api.Issue
}

func (r *testReporter) Report(iss *api.Issue) {
	r.issues = append(r.issues, iss)
}

func (r *testReporter) Debug(msg string, fields ...interface{}) {}

func TestNoLinterDirectiveFilter(t *testing.T) {
	t.Parallel()

	fpath := files.AbsPath("foo.go")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fpath, testSrc, parser.ParseComments)
	assert.NoError(t, err)

	filter := &NoLinterDirectiveFilter{}

	pos := fset.Position(file.Pos())
	filter.AddFile(&api.File{
		Package: &api.Package{
			FSet: fset,
		},
		Position:   &pos,
		CommentMap: ast.NewCommentMap(fset, file, file.Comments),
		ASTFile:    file,
	})

	iss := issue.ToLinterIssue(&api.Issue{
		Position: token.Position{
			Filename: "foo.go",
			Line:     6,
		},
	}, "nolinter")

	assert.True(t, filter.IgnoreIssue(iss))

	iss.Position.Line = 9
	iss.Linter = "foolinter"
	assert.False(t, filter.IgnoreIssue(iss))

	iss.Position.Line = 10
	assert.True(t, filter.IgnoreIssue(iss))

	iss.Linter = "noLinter"
	assert.False(t, filter.IgnoreIssue(iss))

	r := &testReporter{}
	filter.ReportUnnecessaryDirectives(r)
	assert.Len(t, r.issues, 1)
	assert.Equal(t, 14, r.issues[0].Position.Line)
	assert.Equal(t, fpath, r.issues[0].Position.Filename)
}
