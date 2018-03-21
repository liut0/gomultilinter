// Package files contains utils to work with files
package files

import (
	"os"
	"path/filepath"
)

// FileExists returns wether a file exists
func FileExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

// DirExists returns wether a directory exists
func DirExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

// RelPath returns the rel path, errors are ignored
func RelPath(path string) string {
	if !filepath.IsAbs(path) {
		return path
	}

	dir, _ := os.Getwd()
	rel, _ := filepath.Rel(dir, path)
	return rel
}

// AbsPath returns the abs path, errors are ignored
func AbsPath(path string) string {
	abs, _ := filepath.Abs(path)
	return abs
}

// Getwd returns the working directory, errors are ignored
func Getwd() string {
	wd, _ := os.Getwd()
	return wd
}
