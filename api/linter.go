package api

import (
	"context"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/loader"
)

const (
	// LinterFactorySymbolName is the name of the exported Symbol which provides
	// the entry point of a linter plugin
	LinterFactorySymbolName = "LinterFactory"
)

// LinterFactory returns a pointer to a new LinterConfig
type LinterFactory interface {
	NewLinterConfig() LinterConfig
}

// LinterConfig is the struct to which the config gets deserialized
// and which is used to construct new Linter instances
type LinterConfig interface {
	NewLinter() (Linter, error)
}

// Linter base interface.
// concrete impls should implement Linter and one of the ...Linter interfaces below
type Linter interface {
	Name() string
}

// PackageLinter lints packages
type PackageLinter interface {
	Linter
	LintPackage(ctx context.Context, pkg *Package, reporter IssueReporter) error
}

// FileLinter lints files
type FileLinter interface {
	Linter
	LintFile(ctx context.Context, file *File, reporter IssueReporter) error
}

// Package represents a parsed go package
type Package struct {
	PkgInfo *loader.PackageInfo
	FSet    *token.FileSet
}

// File represents a parsed go file
type File struct {
	*Package
	ASTFile    *ast.File
	CommentMap ast.CommentMap
	Position   *token.Position
}
