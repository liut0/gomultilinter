// Package loader downloads, builds, installs and loads linter plugins
package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"plugin"

	"github.com/liut0/gomultilinter/api"
	"github.com/liut0/gomultilinter/config"
	"github.com/liut0/gomultilinter/internal/log"
)

// LoadLinter downloads, installs and loads all the linters specified in the config
func LoadLinter(conf *config.Config) ([]api.Linter, error) {

	if err := os.MkdirAll(conf.LinterInstallDirectory, os.ModePerm); err != nil {
		log.WithFields("install_directory", conf.LinterInstallDirectory).Debug("could not create linter dir")
		return nil, fmt.Errorf("could not create linter directory %s", conf.LinterInstallDirectory)
	}

	linters := make([]api.Linter, 0, len(conf.Linter))
	for _, linterConf := range conf.Linter {
		var (
			linterLibPath string
			err           error
		)

		if linterConf.PluginPath != "" {
			linterLibPath, err = resolvePluginPath(linterConf.PluginPath)
		} else {
			linterLibPath, err = installLinter(linterConf, conf.LinterInstallDirectory, conf.ForceUpdate)
		}

		if err != nil {
			return nil, err
		}

		linter, err := loadLinterPlugin(linterLibPath, linterConf.Config)
		if err != nil {
			return nil, err
		}

		linters = append(linters, linter)
	}

	return linters, nil
}

// loads and initializes the linter plugin
func loadLinterPlugin(libPath string, rawLinterConf []byte) (api.Linter, error) {
	log.WithFields("lib_path", libPath).Debug("loading linter plugin")

	lib, err := plugin.Open(libPath)
	if err != nil {
		log.WithFields("lib_path", libPath, "err", err).Debug("error opening linter lib")
		return nil, fmt.Errorf("error opening linter lib, try -u flag %s: %s", libPath, err)
	}

	symLinterFactory, err := lib.Lookup(api.LinterFactorySymbolName)
	if err != nil {
		log.WithFields("err", err).Debug("linter factory symbol not found")
		return nil, fmt.Errorf("linter factory symbol not found")
	}

	linterFactory, ok := symLinterFactory.(*api.LinterFactory)
	if !ok {
		log.WithFields(
			"type", fmt.Sprintf("%T", symLinterFactory),
			"lib_path", libPath).Debug("linter factory has wrong format")
		return nil, fmt.Errorf("linter factory has wrong format %s: %T", libPath, symLinterFactory)
	}

	lConf := (*linterFactory).NewLinterConfig()
	if rawLinterConf != nil {
		if err := json.Unmarshal(rawLinterConf, lConf); err != nil {
			log.WithFields("lib_path", libPath).Debug("could not unmarshal linter config")
			return nil, fmt.Errorf("linter config unmarshal failed %v", err)
		}
	}

	return lConf.NewLinter()
}
