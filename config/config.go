// Package config manages the configuration of gomultilinter
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/internal/files"
	"github.com/liut0/gomultilinter/internal/log"
	"github.com/sirupsen/logrus"
)

const configFileName = ".gomultilinter.yml"

// Config represents all possible configuration flags for gomultilinter
// see newDefaultConfig for default values
type Config struct {
	// Verbose output
	// only via cli flag
	Verbose bool

	// ForceUpdate enforces rebuild
	// of the linter plugins
	ForceUpdate bool `json:"force_update"`

	// OutputFormat go text/template which is used to print out issues
	// see internal/checker/issue/LinterIssue for available fields
	OutputFormat string `json:"output_format"`

	// MinSeverity for which issues should be printed
	MinSeverity *Severity `json:"min_severity"`

	// Exclude can exclude issues based on their message, name or category
	Exclude *ExcludeConfig `json:"exclude"`

	// LinterInstallDirectory is the dir to which the linter plugins get installed
	LinterInstallDirectory string `json:"linter_install_directory"`

	// Linter which should be used
	Linter []*LinterConfig `json:"linter"`
}

// ExcludeConfig excludes Issues by matching
// corresponding regular expressions
type ExcludeConfig struct {
	// if true unnecessary nolint comment directives do not result in an issue
	UnnecessaryNoLintDirectives bool `json:"unnecessary_no_lint_directives"`

	// Tests if true test files are ignored
	Tests bool `json:"tests"`

	// File-/Packagenames which should be excluded
	Names MultiRegex `json:"names"`

	// Linter messages which should be excluded
	Messages MultiRegex `json:"messages"`

	// Linter cagtegories which should be excluded
	Categories MultiRegex `json:"categories"`
}

// LinterConfig represents a Linter which should be used
// Package or PluginPath needs to be provided
type LinterConfig struct {
	// Package of the gomultilinter plugin
	Package string `json:"package"`

	// PluginPath is the path to the .so file of the gomultilinter plugin
	PluginPath string `json:"plugin_path"`

	// Config is the Configuration for the concrete linter
	Config json.RawMessage `json:"config"`
}

func newDefaultConfig() *Config {
	return &Config{
		MinSeverity:            &Severity{Severity: api.SeverityInfo},
		LinterInstallDirectory: os.ExpandEnv("$GOPATH/pkg/gomultilinter/linter"),
		OutputFormat:           "{{.Path}}:{{.Line}}:{{if .Col}}{{.Col}}{{end}}:{{.Severity}}:{{.Category}}: {{.Message}} ({{.Linter}})",
		Exclude:                new(ExcludeConfig),
	}
}

// SetVerbose initializes the logger with the corresponding severity
func SetVerbose(verbose bool) {
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

// ReadConfig reads the config file at the specified path
// if the specified path is empty it searches from the current
// directory upwards in the fs for a file named .gomultilinter.yml
// verbose, and forceUpdate are cli flags which can be provided and
// override the flags from the configfile
func ReadConfig(path string, verbose, forceUpdate bool) (*Config, error) {
	if path == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.WithFields("err", err).Debug("config: failed to read working directory")
			return nil, errors.New("reading config failed")
		}
		path = findConfigFile(wd)
	} else if !files.FileExists(path) {
		return nil, fmt.Errorf("config file at %v not found", path)
	}

	log.WithFields("path", path).Debug("using config file")

	conf := newDefaultConfig()

	confFileContent, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithFields("err", err).Fatal("could not read config file")
	}

	if err := yaml.Unmarshal(confFileContent, conf); err != nil {
		log.WithFields("err", err).Fatal("could not read config file")
	}

	// overwrite cli flags
	conf.Verbose = verbose
	if forceUpdate {
		conf.ForceUpdate = forceUpdate
	}

	return conf, nil
}

func findConfigFile(path string) string {
	file := filepath.Join(path, configFileName)
	if !files.FileExists(file) {
		parent := filepath.Dir(path)
		if parent == path {
			return ""
		}
		return findConfigFile(parent)
	}
	return file
}
