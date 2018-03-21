package loader

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/liut0/gomultilinter/config"
	"github.com/liut0/gomultilinter/internal/files"
	"github.com/liut0/gomultilinter/internal/log"
)

// installLinter downloads a package (if not available locally)
// and builds (if not yet builded or forceBuild is set) the *.so file to installDir
func installLinter(linterConf *config.LinterConfig, installDir string, forceBuild bool) (string, error) {
	pkgImportPath, foundLocally := resolveImportPath(linterConf.Package)

	if !foundLocally {
		if err := downloadPkg(pkgImportPath); err != nil {
			return "", err
		}
	}

	return buildPlugin(pkgImportPath, installDir, forceBuild)
}

func buildPlugin(pkg, installDir string, forceBuild bool) (string, error) {

	libPath := filepath.Join(installDir, pkg) + ".so"
	libDir := filepath.Dir(libPath)

	// TODO improve decision wether a new build is required (plugin load
	// fails if api pkg changes) check go.buildID?
	if !forceBuild && files.FileExists(libPath) {
		log.WithFields("lib_path", libPath).Debug("is up to date")
		return libPath, nil
	}

	if err := os.MkdirAll(libDir, os.ModePerm); err != nil {
		return "", err
	}

	log.WithFields("linter_pkg", pkg).Debug("go build")
	if err := execGoCommand("build", "--buildmode", "plugin", "-v", "-o", libPath, pkg); err != nil {
		return "", err
	}

	return libPath, nil
}

func downloadPkg(pkg string) error {
	log.WithFields("linter_pkg", pkg).Debug("go get")
	return execGoCommand("get", "-v", "-d", pkg)
}

func execGoCommand(args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func resolveImportPath(pkg string) (string, bool) {
	importPkg, err := build.Import(pkg, files.Getwd(), 0)
	if err != nil {
		return pkg, false
	}
	return importPkg.ImportPath, true
}

// resolvePluginPath expands the path and checks wether a file
// at the given location exists
func resolvePluginPath(pluginPath string) (string, error) {
	pluginPath = os.ExpandEnv(pluginPath)
	if files.FileExists(pluginPath) {
		return pluginPath, nil
	}

	log.WithFields("plugin_path", pluginPath).Debug("plugin not found")
	return "", fmt.Errorf("plugin not found %s", pluginPath)
}
