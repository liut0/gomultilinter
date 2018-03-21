// Package filter provides filters to filter out excluded issues
package filter

import (
	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/config"
	"github.com/liut0/gomultilinter/internal/checker/issue"
)

// IssueFilter can filter out issues
type IssueFilter interface {
	IgnoreIssue(issue *issue.LinterIssue) bool
}

// IssueFilterFunc wraps a func to satisfy the IssueFilter interface
type IssueFilterFunc func(issue *issue.LinterIssue) bool

// IgnoreIssue just wraps the provided func
func (f IssueFilterFunc) IgnoreIssue(issue *issue.LinterIssue) bool {
	return f(issue)
}

// ChainFilter chains multiple filters
func ChainFilter(filters ...IssueFilter) IssueFilter {
	return IssueFilterFunc(func(issue *issue.LinterIssue) bool {
		for _, filter := range filters {
			if filter.IgnoreIssue(issue) {
				return true
			}
		}
		return false
	})
}

// SeverityFilter returns an IssueFilter which filters out issues
// with a lower severity than minSeverity
func SeverityFilter(minSeverity api.Severity) IssueFilter {
	return IssueFilterFunc(func(issue *issue.LinterIssue) bool {
		return minSeverity > issue.Severity
	})
}

// MessageFilter returns an IssueFilter which filters out issues
// with a message matching any of the provided regular expressions
func MessageFilter(exclude config.MultiRegex) IssueFilter {
	return IssueFilterFunc(func(issue *issue.LinterIssue) bool {
		return exclude.MatchesAny(issue.Message)
	})
}

// CategoryFilter returns an IssueFilter which filters out issues
// with a category matching any of the provided regular expressions
func CategoryFilter(exclude config.MultiRegex) IssueFilter {
	return IssueFilterFunc(func(issue *issue.LinterIssue) bool {
		return exclude.MatchesAny(issue.Category)
	})
}

// FilenameFilter returns an IssueFilter which filters out issues
// with a filename matching any of the provided regular expressions
func FilenameFilter(exclude config.MultiRegex) IssueFilter {
	return IssueFilterFunc(func(issue *issue.LinterIssue) bool {
		return exclude.MatchesAny(issue.Path.Abs)
	})
}
