package imports

import (
	"errors"
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/liut0/gomultilinter/config"
	"github.com/liut0/gomultilinter/internal/files"
	"github.com/liut0/gomultilinter/internal/log"
)

const (
	recursiveSuffix = "/..."
)

var (
	ignoredDirs = config.MultiRegex{
		&config.Regex{Regexp: regexp.MustCompile(`\.git`)},
		&config.Regex{Regexp: regexp.MustCompile(`vendor`)},
	}
)

// ResolvePaths resolves paths to go packages
// paths can either be directories, packages or go files
//
// if a directory path or a package has a '/...' suffix all
// subpackages/-directories are also included
//
// if paths is empty the current directory is used including all
// subdirectories
func ResolvePaths(paths ...string) ([]string, error) {
	if len(paths) == 0 {
		paths = []string{"." + recursiveSuffix}
	}

	imports := make([]string, 0, len(paths))

	var file, dir, pkg int
	for _, path := range paths {
		rec := strings.HasSuffix(path, recursiveSuffix)
		cPath := path
		if rec {
			cPath = strings.TrimSuffix(cPath, recursiveSuffix)
		}

		switch {
		case files.DirExists(cPath):
			dir = 1
			pkgs, err := resolveDir(cPath, rec)
			if err != nil {
				return nil, err
			}
			imports = append(imports, pkgs...)
		case files.FileExists(path):
			file = 1
			imports = append(imports, path)
		default:
			pkg = 1
			if rec {
				pkgs, err := resolveRecursivePkg(cPath)
				if err != nil {
					return nil, err
				}

				imports = append(imports, pkgs...)
			} else {
				imports = append(imports, path)
			}
		}
	}

	if file+dir+pkg != 1 {
		log.WithFields("file", file, "dir", dir, "pkg", pkg).Debug("multiple target types")
		return nil, errors.New("multiple target types, ensure flags are before targets/paths")
	}

	return imports, nil
}

func getRecursiveSubDirs(dir string) []string {
	var paths []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !ignoredDirs.MatchesAny(path) {
			paths = append(paths, path)
		}
		return nil
	})

	return paths
}

func expandRecursivePath(path string, recursive bool) []string {
	if !recursive {
		return []string{path}
	}

	return getRecursiveSubDirs(path)
}

func resolveDir(path string, recursive bool) ([]string, error) {
	paths := expandRecursivePath(files.AbsPath(path), recursive)

	pkgs := make([]string, 0, len(paths))
	for _, path := range paths {

		pkg, err := build.ImportDir(path, 0)
		if err != nil {
			if _, isNoGo := err.(*build.NoGoError); isNoGo {
				continue
			}

			log.WithFields("err", err, "path", path).Debug("error importing directory")
			return nil, err
		}

		pkgs = append(pkgs, pkg.ImportPath)
	}

	return pkgs, nil
}

func resolveRecursivePkg(pkgImportPath string) ([]string, error) {
	pkg, err := build.Import(pkgImportPath, ".", build.FindOnly)
	if err != nil {
		log.WithFields("err", err, "pkg", pkgImportPath).Debug("error importing pkgImportPath")
		return nil, err
	}

	pkgs, err := resolveDir(pkg.Dir, true)
	if err != nil {
		return nil, err
	}

	return pkgs, nil
}
