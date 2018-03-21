package filter

import (
	"regexp"
	"testing"

	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/config"
	"github.com/liut0/gomultilinter/internal/checker/issue"
	"github.com/stretchr/testify/assert"
)

func TestChainFilter(t *testing.T) {
	t.Parallel()

	f := ChainFilter(
		IssueFilterFunc(func(issue *issue.LinterIssue) bool {
			return issue.Message == "1"
		}),
		IssueFilterFunc(func(issue *issue.LinterIssue) bool {
			return issue.Message == "2"
		}),
	)

	assert.True(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Message: "1"}}))
	assert.True(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Message: "2"}}))
	assert.False(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Message: "3"}}))
}

func TestSeverityFilter(t *testing.T) {
	t.Parallel()

	f := SeverityFilter(api.SeverityWarning)

	assert.True(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Severity: api.SeverityInfo}}))
	assert.False(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Severity: api.SeverityWarning}}))
	assert.False(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Severity: api.SeverityError}}))
}

func TestMessageFilter(t *testing.T) {
	t.Parallel()

	f := MessageFilter(config.MultiRegex{
		&config.Regex{Regexp: regexp.MustCompile("foo")},
		&config.Regex{Regexp: regexp.MustCompile("bar")},
	})

	assert.True(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Message: "foo"}}))
	assert.True(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Message: "bar"}}))
	assert.False(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Message: "xy"}}))
}

func TestCategoryFilter(t *testing.T) {
	t.Parallel()

	f := CategoryFilter(config.MultiRegex{
		&config.Regex{Regexp: regexp.MustCompile("foo")},
		&config.Regex{Regexp: regexp.MustCompile("bar")},
	})

	assert.True(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Category: "foo"}}))
	assert.True(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Category: "bar"}}))
	assert.False(t, f.IgnoreIssue(&issue.LinterIssue{Issue: &api.Issue{Category: "xy"}}))
}

func TestFilenameFilter(t *testing.T) {
	t.Parallel()

	f := FilenameFilter(config.MultiRegex{
		&config.Regex{Regexp: regexp.MustCompile("foo")},
		&config.Regex{Regexp: regexp.MustCompile("bar")},
	})

	assert.True(t, f.IgnoreIssue(&issue.LinterIssue{Path: issue.Path{Abs: "/foo/xy", Rel: "./xy"}}))
	assert.True(t, f.IgnoreIssue(&issue.LinterIssue{Path: issue.Path{Abs: "/x/y/z/bar/hello", Rel: "../../z/bar/hello"}}))
	assert.False(t, f.IgnoreIssue(&issue.LinterIssue{Path: issue.Path{Abs: "/x/y/z/hello", Rel: "../../z/hello"}}))
}
