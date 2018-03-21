package main

import (
	"go/token"
	"os"
	"reflect"
	"testing"

	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/config"
	"github.com/liut0/gomultilinter/internal/checker"
	"github.com/liut0/gomultilinter/internal/checker/issue"
	"github.com/liut0/gomultilinter/internal/files"
	"github.com/liut0/gomultilinter/internal/loader"
	"github.com/stretchr/testify/assert"
)

func TestMainCMD(t *testing.T) {
	assert.Equal(t, 2, mainCMD(&flags{
		configFile:  os.ExpandEnv("$GOPATH/src/github.com/liut0/gomultilinter/test/data/.gomultilinter.yml"),
		forceUpdate: true}))

	assert.Equal(t, 0, mainCMD(&flags{
		configFile:   os.ExpandEnv("$GOPATH/src/github.com/liut0/gomultilinter/test/data/.gomultilinter.yml"),
		forceUpdate:  true,
		noExitStatus: true,
	}))
}

func TestMainAll(t *testing.T) {
	conf, err := config.ReadConfig(
		os.ExpandEnv("$GOPATH/src/github.com/liut0/gomultilinter/test/data/.gomultilinter.yml"),
		false,
		true)
	assert.NoError(t, err)

	linter, err := loader.LoadLinter(conf)
	assert.NoError(t, err)

	ckr, err := checker.NewChecker(conf, linter)
	assert.NoError(t, err)

	err = ckr.Load("github.com/liut0/gomultilinter/test/data")
	assert.NoError(t, err)

	issues := ckr.Run()

	expectIssues(t, []*issue.LinterIssue{
		issue.ToLinterIssue(&api.Issue{
			Category: "comments",
			Message:  "exported function TestIssue2 should have comment or be unexported",
			Severity: api.SeverityWarning,
			Position: token.Position{
				Filename: files.AbsPath("test/data/issue.go"),
				Line:     7,
				Column:   1,
				Offset:   46,
			},
		}, "golint"),
		issue.ToLinterIssue(&api.Issue{
			Category: "unnecessary-nolinter-directive",
			Message:  "unnecessary nolinter directive detected",
			Severity: api.SeverityWarning,
			Position: token.Position{
				Filename: files.AbsPath("test/data/issue.go"),
				Line:     16,
				Column:   1,
				Offset:   134,
			},
		}, "gomultilinter"),
		issue.ToLinterIssue(&api.Issue{
			Category: "cyclo",
			Message:  "cyclomatic complexity [8/1] of function testIssue5 is high",
			Severity: api.SeverityWarning,
			Position: token.Position{
				Filename: files.AbsPath("test/data/issue.go"),
				Line:     19,
				Column:   1,
				Offset:   157,
			},
		}, "GoCyclo"),
		issue.ToLinterIssue(&api.Issue{
			Category: "test",
			Message:  "testmsg",
			Severity: api.SeverityWarning,
			Position: token.Position{
				Filename: files.AbsPath("test/data/issue.go"),
				Line:     4,
				Column:   2,
				Offset:   34,
			},
		}, "testlinter"),
		issue.ToLinterIssue(&api.Issue{
			Category: "linter-panic",
			Message:  "linter did panic: testmsg",
			Severity: api.SeverityError,
		}, "testlinter"),
		issue.ToLinterIssue(&api.Issue{
			Category: "linter-error",
			Message:  "linter returned error: testmsg",
			Severity: api.SeverityError,
		}, "testlinter"),
	}, issues)
}

func expectIssues(t *testing.T, expected []*issue.LinterIssue, actual []*issue.LinterIssue) {
	assert.Len(t, actual, len(expected))
	for _, iss := range expected {
		var found bool
		for _, actual := range actual {
			if reflect.DeepEqual(iss, actual) {
				found = true
				break
			}
		}

		if !found {
			assert.Failf(t, "expected issue not found", "%v", iss)
		}
	}
}
