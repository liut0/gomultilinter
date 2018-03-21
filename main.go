package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/liut0/gomultilinter/config"
	"github.com/liut0/gomultilinter/internal/checker"
	"github.com/liut0/gomultilinter/internal/loader"
	"github.com/liut0/gomultilinter/internal/log"
)

const (
	exitSuccess = 0

	// 2 becaus log.Fatal() uses 1
	exitIssues = 2
)

type flags struct {
	configFile   string
	verbose      bool
	forceUpdate  bool
	installOnly  bool
	noExitStatus bool
}

func usage() {
	fmt.Fprintf(os.Stderr,
		`usage: %s [flags] <targets>

targets:
    none         current directory including all sub-directoreis, same as './...'
    packages     where a '/...' suffix includes all sub-packages
    directories  where a '/...' suffix includes all sub-directories
    files        all must belong to a single package

flags:
`, os.Args[0])

	flag.PrintDefaults()
}

func main() {
	cliFlags := &flags{}
	flag.Usage = usage
	flag.StringVar(&cliFlags.configFile, "config", "", "set path of the config file")
	flag.BoolVar(&cliFlags.verbose, "v", false, "verbose output")
	flag.BoolVar(&cliFlags.forceUpdate, "u", false, "force update/rebuild of linters")
	flag.BoolVar(&cliFlags.installOnly, "install-only", false, "build/install/validate plugins only, do not lint")
	flag.BoolVar(&cliFlags.noExitStatus, "no-exit-status", false, "sets exit status only to non 0 if an underlying error occurs")
	flag.Parse()

	os.Exit(mainCMD(cliFlags))
}

func mainCMD(cliFlags *flags) int {
	metrics := newMetrics()
	metricsMain := metrics.newEntry("main")
	defer metrics.log()
	defer metricsMain.done()

	config.SetVerbose(cliFlags.verbose)

	metricsConf := metrics.newEntry("conf")
	conf, err := config.ReadConfig(cliFlags.configFile, cliFlags.verbose, cliFlags.forceUpdate)
	if err != nil {
		log.WithFields("err", err).Fatal()
	}
	metricsConf.done()

	metricsLoadPlugins := metrics.newEntry("load_plugins")
	linter, err := loader.LoadLinter(conf)
	if err != nil {
		log.WithFields("err", err).Fatal()
	}
	metricsLoadPlugins.done()
	if cliFlags.installOnly {
		return exitSuccess
	}

	metricsLoadChecker := metrics.newEntry("load_checker")
	ckr, err := checker.NewChecker(conf, linter)
	if err != nil {
		log.WithFields("err", err).Fatal()
	}

	if err := ckr.Load(flag.Args()...); err != nil {
		log.WithFields("err", err).Fatal()
	}
	metricsLoadChecker.done()

	metricsLinters := metrics.newEntry("linters")
	issues := ckr.Run()
	metricsLinters.done()

	issuesCount := len(issues)

	log.WithFields("issues_count", issuesCount).Debug("done")

	if cliFlags.noExitStatus || issuesCount == 0 {
		return exitSuccess
	}

	return exitIssues
}
