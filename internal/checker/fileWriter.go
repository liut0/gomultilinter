package checker

import (
	"fmt"
	"os"
	"text/template"

	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/internal/checker/issue"
	"github.com/liut0/gomultilinter/internal/log"
)

// IssueWriter writes an issue to a target
type IssueWriter interface {
	Write(issue *issue.LinterIssue)
}

// FileWriter implements the IssueWriter interface
// and writes the issues serialized by the outtemplate
// to the provided files
type FileWriter struct {
	outTemplate *template.Template
	out         map[api.Severity]*os.File
}

func newConsoleWriter(tmplStr string) (*FileWriter, error) {
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		log.WithFields("err", err).Debug("could not parse output template")
		return nil, fmt.Errorf("could not parse output template %v", err)
	}

	return &FileWriter{
		outTemplate: tmpl,
		out: map[api.Severity]*os.File{
			api.SeverityInfo:    os.Stdout,
			api.SeverityWarning: os.Stderr,
			api.SeverityError:   os.Stderr,
		},
	}, nil
}

func (w *FileWriter) Write(issue *issue.LinterIssue) {
	out, ok := w.out[issue.Severity]
	if ok {
		if err := w.outTemplate.Execute(out, issue); err != nil {
			log.WithFields("err", err).Error("output template execution failed")
		}
		fmt.Fprintln(out)
	}
}
