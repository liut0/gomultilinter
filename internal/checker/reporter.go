package checker

import (
	"sync"

	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/internal/checker/filter"
	"github.com/liut0/gomultilinter/internal/checker/issue"
	"github.com/liut0/gomultilinter/internal/log"
)

// IssueReporter implements the api.Reporter interface
// and filters/collects issues and passes them to the writer
type IssueReporter struct {
	filter filter.IssueFilter

	allIssuesLock sync.Mutex
	allIssues     []*issue.LinterIssue

	issueWriter IssueWriter
}

// IssueReporterEntry is a concrete reporter of a linter's invocation
type IssueReporterEntry struct {
	*IssueReporter
	linter string
}

func (r *IssueReporter) entry(linter string) *IssueReporterEntry {
	return &IssueReporterEntry{
		IssueReporter: r,
		linter:        linter,
	}
}

// Debug prints debug messages
func (r *IssueReporterEntry) Debug(msg string, fields ...interface{}) {
	fields = append(fields, "linter", r.linter)
	log.WithFields(fields...).Debug(msg)
}

// Report checks if an issue gets filtered
// if so the issue is ignored
// otherwise it adds the issue to the list of all issues and
// passes it to the writer
func (r *IssueReporterEntry) Report(iss *api.Issue) {
	linterIssue := issue.ToLinterIssue(iss, r.linter)

	if r.filter.IgnoreIssue(linterIssue) {
		return
	}

	r.addIssue(linterIssue)

	r.issueWriter.Write(linterIssue)
}

func (r *IssueReporterEntry) addIssue(issue *issue.LinterIssue) {
	r.allIssuesLock.Lock()
	defer r.allIssuesLock.Unlock()

	r.allIssues = append(r.allIssues, issue)
}
