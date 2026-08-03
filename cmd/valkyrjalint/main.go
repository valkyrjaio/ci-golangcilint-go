/*
 * This file is part of the Valkyrja golangci-lint package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Command valkyrjalint writes the shared Valkyrja Go CI configuration, and it
// writes the copyright header into a Go file.
//
// A repository runs the command through `go tool`, so the version is pinned in
// the repository's own tool module and never mixes with its dependencies:
//
//	go tool -modfile=.github/ci/lint/go.mod valkyrjalint config > .golangci.yml
//	go tool -modfile=.github/ci/lint/go.mod valkyrjalint header -package 'Valkyrja Framework' -w .
package main

import (
	"os"

	"github.com/valkyrjaio/ci-golangcilint-go/v26/golangcilint/cli"
)

// exit ends the process. A test replaces it, so this file carries a test of its
// own rather than an exclusion from the coverage report.
var exit = os.Exit

func main() {
	exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
