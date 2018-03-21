// Package checker can load files/packages and run linters on them
package checker

import (
	"fmt"
	"go/parser"
	"go/token"
	"go/types"

	"context"

	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/config"
	"github.com/liut0/gomultilinter/internal/checker/filter"
	"github.com/liut0/gomultilinter/internal/checker/imports"
	"github.com/liut0/gomultilinter/internal/checker/issue"
	"github.com/liut0/gomultilinter/internal/log"
	"golang.org/x/tools/go/loader"
)

const (
	selfLinterName = "gomultilinter"
)

// Checker is the coordinator/executor of the linting process
type Checker struct {
	excludeUnnecessaryNoLintDirectives bool
	excludeTests                       bool
	excludeNames                       config.MultiRegex

	fileLinter map[string]api.FileLinter
	pkgLinter  map[string]api.PackageLinter

	issueReporter           *IssueReporter
	selfIssueReporter       api.IssueReporter
	noLinterDirectiveFilter *filter.NoLinterDirectiveFilter

	pkgs []*loader.PackageInfo
	fset *token.FileSet

	ctx context.Context
}

// NewChecker constructs a new checker according to the provided arguments
func NewChecker(conf *config.Config, linter []api.Linter) (*Checker, error) {
	issueWriter, err := newConsoleWriter(conf.OutputFormat)
	if err != nil {
		return nil, err
	}

	noLinterDirectiveFilter := &filter.NoLinterDirectiveFilter{}

	reporter := &IssueReporter{
		issueWriter: issueWriter,

		filter: filter.ChainFilter(
			filter.SeverityFilter(conf.MinSeverity.Severity),
			filter.CategoryFilter(conf.Exclude.Categories),
			// filter names again (pkglinters cant filter filenames before linting)
			filter.FilenameFilter(conf.Exclude.Names),
			filter.MessageFilter(conf.Exclude.Messages),
			noLinterDirectiveFilter),
	}

	c := &Checker{
		excludeUnnecessaryNoLintDirectives: conf.Exclude.UnnecessaryNoLintDirectives,
		excludeTests:                       conf.Exclude.Tests,
		excludeNames:                       conf.Exclude.Names,

		fileLinter: map[string]api.FileLinter{},
		pkgLinter:  map[string]api.PackageLinter{},

		issueReporter:           reporter,
		selfIssueReporter:       reporter.entry(selfLinterName),
		noLinterDirectiveFilter: noLinterDirectiveFilter,

		ctx: context.Background(),
	}

	for _, l := range linter {
		switch lT := l.(type) {
		case api.FileLinter:
			c.fileLinter[lT.Name()] = lT
		case api.PackageLinter:
			c.pkgLinter[lT.Name()] = lT
		default:
			log.WithFields("linter", l.Name()).Debug("unsupported linter")
			return nil, fmt.Errorf("unsupported linter %v", l.Name())
		}
	}

	return c, nil
}

// Load loads/parses the specified paths
// see imports.ResolvePaths how paths are resolved
func (c *Checker) Load(paths ...string) error {
	log.WithFields("pahts", paths).Debug("loading paths")
	importPaths, err := imports.ResolvePaths(paths...)
	if err != nil {
		return err
	}

	program, err := c.load(importPaths)
	if err != nil {
		log.WithFields("err", err).Debug("could not load pkgs")
		return fmt.Errorf("could not load pkgs %v", err)
	}

	c.pkgs = program.InitialPackages()
	c.fset = program.Fset
	return nil
}

// Run runs the checker on the loaded paths
func (c *Checker) Run() []*issue.LinterIssue {
	log.Debug("running linters")

	c.walkPkgs(c.pkgs, c.fset)

	if !c.excludeUnnecessaryNoLintDirectives {
		c.noLinterDirectiveFilter.ReportUnnecessaryDirectives(c.selfIssueReporter)
	}

	return c.issueReporter.allIssues
}

func (c *Checker) load(paths []string) (*loader.Program, error) {
	loadCfg := loader.Config{
		TypeChecker: types.Config{
			// ignore type checker errors
			Error: func(error) {},
		},
		AllowErrors: true,
		ParserMode:  parser.ParseComments,
	}

	if r, err := loadCfg.FromArgs(paths, !c.excludeTests); err != nil {
		return nil, err
	} else if len(r) > 0 {
		log.WithFields("args", r).Debug("invalid arguments")
		return nil, fmt.Errorf("invalid arguments %v", r)
	}

	return loadCfg.Load()
}
