<p align="center">
    <h3 align="center">GoMultiLinter</h3>
    <p align="center">
        <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
        <a href="https://travis-ci.org/liut0/gomultilinter"><img alt="Travis" src="https://img.shields.io/travis/liut0/gomultilinter/master.svg?style=flat-square"></a>
        <a href="https://goreportcard.com/report/github.com/liut0/gomultilinter"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/liut0/gomultilinter?style=flat-square"></a>
    </p>
</p>

---

Run Go lint tools via a plugin api. Parse AST once, lint it by all the linters.

Since [Plugins](https://golang.org/pkg/plugin/) are used under the hood gomultilinter currently only supports linux and macOS.

- [Quickstart](#quickstart)
- [Installation](#installation)
- [Editor integration](#editor-integration)
- [Configuration](#configuration)
    - [Example configuration file](#example-configuration-file)
- [Comment directives](#comment-directives)
- [Linters](#linters)
    - [Available Linters](#available-linters)
    - [Custom Linters](#custom-linters)
    - [Linter Vendoring](#linter-vendoring)
- [Exit status](#exit-status)
- [Feedback](#feedback)

## Quickstart

- [Install](#installation) `gomultilinter`
- [Create a config file](#configuration)
- Run `gomultilinter`

## Installation

go get from HEAD: `go get -u github.com/liut0/gomultilinter`

Ensure the linter's are built with the excact same gomultilinter source as the gomultilinter binary.

## Editor integration

- Intellij: Use the `filewatchers` with the `gomultilinter` template and override the settings below:
    - Program: `gomultilinter`
    - Arguments: `$FilePath$`
    - Working directory: `$FileDir$`

## Configuration

gomultilinter is configured via a yaml config file. Either the name of the configuration file has to be `.gomultilinter.yml` and it must be placed in the working directory or any parent directory or the location of the config file can be passed by the cli flag:
`-config=<file>`. The format of this file is determined by
the `Config` struct in [config.go](https://github.com/liut0/gomultilinter/blob/master/config/config.go).

### Example configuration file

```yaml
exclude:
  tests: true
  names:
    - '_mock\\.go'
  categories:
    - 'comments'

linter:
  - package: 'github.com/liut0/gomultilinter-golint/gomultilinter'
    config:
      - max_cyclo: 10
  - plugin_path: '~/myGoMultilinterPlugin.so'
    config:
      - foo: 'bar'
```

## Comment directives

gomultilinter supports suppression of linter messages via comment directives. The
form of the directive is:

```go
// nolint[: <linter>[, <linter>, ...]]
```

## Linters

### Available Linters

Currently the available linters are forks which implement the required interfaces 
with preferably minimal (but sometimes dirty, for easier merging) changes to the main repo.
Maybe someday the maintainers of the linters will implement the interface in the main repository and the forks are not needed anymore.
More linters will be added...

- golint: `github.com/liut0/gomultilinter-golint/gomultilinter`
    - Fork of [github.com/golang/lint](https://github.com/golang/lint)
    - Config: `min_confidence`
- GoCyclo: `github.com/liut0/gomultilinter-gocyclo`
    - Fork of [github.com/fzipp/gocyclo](https://github.com/fzipp/gocyclo)
    - Config: `max_cyclo`
- errcheck: `github.com/liut0/gomultilinter-errcheck/gomultilinter`
    - Fork of [github.com/kisielk/errcheck](https://github.com/kisielk/errcheck)
    - Config: `blank`, `asserts`, `exclude`
- maligned: `github.com/liut0/gomultilinter-maligned`
    - Fork of [github.com/mdempsky/maligned](https://github.com/mdempsky/maligned)
- deadcode: `github.com/liut0/gomultilinter-deadcode`
    - Fork of [github.com/tsenart/deadcode](https://github.com/tsenart/deadcode)
- preventusage: `github.com/liut0/gomultilinter-commonlinters/preventusage`
    - Config: [see gomultilinter-commonlinters/preventusage](https://github.com/liut0/gomultilinter-commonlinters)
- dep: `github.com/liut0/gomultilinter-commonlinters/dep`
     - Config: [see gomultilinter-commonlinters/dep](https://github.com/liut0/gomultilinter-commonlinters)
- licenses: `github.com/liut0/gomultilinter-commonlinters/licenses`
     - Config: [see gomultilinter-commonlinters/licenses](https://github.com/liut0/gomultilinter-commonlinters)
        
### Custom Linters

Implementation

- Implement the interfaces described in [api/linter.go](https://github.com/liut0/gomultilinter/blob/master/api/linter.go). 
- In a package named `main` hold a variable called `LinterFactory` of the type `github.com/liut0/api/linter/LinterFactory`.

The custom linter can be added to gomultilinter in two ways:

- via the `Package` configuration directive of the config file (the package gets built by gomultilinter with `buildmode=plugin`)
- via the `PluginPath` configuration directive of the config file (path to the prebuilt `.so` plugin file which gets picked up by gomultilinter)

### Linter Vendoring

Linter packages get resolved from the working directory the same way go does. If a linter pkg exists in the vendor dir it's preferred.
Therefore you can vendor the linter's how you would vendor any other package.

## Exit status

| Value | Meaning |
| - | - |
| 0 | Succeed :) |
| 1 | An underlying error occurred |
| 2 | Issues occurred. Unless `no-exit-status` cli flag is set. |

## Feedback

Feedback is greatly appreciated. If you have any questions, please don't hesitate to create a issue.

inspired by [gometalinter](https://github.com/alecthomas/gometalinter)