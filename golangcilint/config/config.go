/*
 * This file is part of the Valkyrja golangci-lint package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package config holds the golangci-lint configuration that every Valkyrja Go
// repository runs.
//
// golangci-lint reads one YAML file and has no way to extend another, so a
// repository cannot import this configuration the way a PHP repository imports
// `Valkyrja\Fixer\Rules`. This package embeds the file instead, and the
// `valkyrjalint config` command writes it. A repository commits the result and
// checks it, the same way `go mod tidy -diff` checks a generated file.
package config

import _ "embed"

// golangci holds the configuration file this package ships.
//
//go:embed golangci.yml
var golangci string

// Get returns the golangci-lint configuration that every Valkyrja Go repository
// runs. The result ends with a newline.
func Get() string {
	return golangci
}
