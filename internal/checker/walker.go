package checker

import (
	"go/ast"
	"go/token"
	"regexp"

	"github.com/liut0/gomultilinter/api"
	"golang.org/x/tools/go/loader"
)

var (
	generatedFileRgx = regexp.MustCompile("(?si).*code generated.*do not edit.*")
)

func (c *Checker) walkPkgs(pkgInfos []*loader.PackageInfo, fset *token.FileSet) {
	for _, pkgInfo := range pkgInfos {
		c.walkPkg(pkgInfo, fset)
	}
}

func (c *Checker) walkPkg(pkgInfo *loader.PackageInfo, fset *token.FileSet) {

	pkg := &api.Package{
		PkgInfo: pkgInfo,
		FSet:    fset,
	}

	if c.ignorePkg(pkg) {
		return
	}

	files := make([]*api.File, 0, len(pkgInfo.Files))

	for _, astFile := range pkgInfo.Files {
		fpos := pkg.FSet.Position(astFile.Pos())
		file := &api.File{
			Package:    pkg,
			ASTFile:    astFile,
			Position:   &fpos,
			CommentMap: ast.NewCommentMap(pkg.FSet, astFile, astFile.Comments),
		}

		if !c.ignoreFile(file) {
			c.noLinterDirectiveFilter.AddFile(file)
			files = append(files, file)
		}
	}

	c.lintPkg(pkg)

	for _, f := range files {
		c.lintFile(f)
	}
}

func (c *Checker) ignorePkg(pkg *api.Package) bool {
	return c.excludeNames.MatchesAny(pkg.PkgInfo.Pkg.Path())
}

func (c *Checker) ignoreFile(file *api.File) bool {

	if c.excludeNames.MatchesAny(file.Position.Filename) {
		return true
	}

	for _, cmnt := range file.ASTFile.Comments {
		if generatedFileRgx.MatchString(cmnt.Text()) {
			return true
		}
	}
	return false
}
