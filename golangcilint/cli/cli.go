/*
 * This file is part of the Valkyrja golangci-lint package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package cli holds the `valkyrjalint` commands.
//
// The package holds the whole command, and `cmd/valkyrjalint` holds one call to
// Run. A test therefore drives every command through an ordinary function call
// and reads the output from a buffer.
package cli

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/valkyrjaio/ci-golangcilint-go/v26/golangcilint/config"
	"github.com/valkyrjaio/ci-golangcilint-go/v26/golangcilint/header"
)

// usage describes every command.
const usage = `valkyrjalint writes the shared Valkyrja Go CI configuration.

Usage:
  valkyrjalint config
        Write the golangci-lint configuration to stdout.

  valkyrjalint header -package NAME [-w] PATH...
        Report each Go file whose copyright header is missing or wrong.
        -w rewrites each such file instead of reporting it.
`

// fileMode is the permission a rewritten file keeps.
const fileMode = 0o600

// skipped names each directory the header walk never descends into. A file under
// one of these is not this repository's source.
var skipped = map[string]bool{
	".git":       true,
	".worktrees": true,
	"vendor":     true,
	"testdata":   true,
}

// writef puts a message on a stream.
//
// A diagnostic that fails to write has nowhere left to report the failure, so
// this function drops the error rather than returning one no caller could act
// on.
func writef(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

// Run performs one command and returns the process status.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		writef(stderr, "%s", usage)

		return 2
	}

	switch args[0] {
	case "config":
		return runConfig(args[1:], stdout, stderr)
	case "header":
		return runHeader(args[1:], stdout, stderr)
	default:
		writef(stderr, "unknown command %q\n\n%s", args[0], usage)

		return 2
	}
}

// runConfig writes the shared golangci-lint configuration.
func runConfig(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		writef(stderr, "config takes no argument, got %q\n", args[0])

		return 2
	}

	writef(stdout, "%s", config.Get())

	return 0
}

// runHeader reports or rewrites the copyright header of each Go file it reads.
func runHeader(args []string, stdout, stderr io.Writer) int {
	set := flag.NewFlagSet("header", flag.ContinueOnError)
	set.SetOutput(stderr)

	packageIdentifier := set.String("package", "", "the package identifier the header names")
	rewrite := set.Bool("w", false, "rewrite each file instead of reporting it")

	err := set.Parse(args)
	if err != nil {
		return 2
	}

	if *packageIdentifier == "" {
		writef(stderr, "header needs -package. COPYRIGHT_HEADER.md names the value for each repository.\n")

		return 2
	}

	paths := set.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	files, err := collect(paths)
	if err != nil {
		writef(stderr, "%v\n", err)

		return 1
	}

	return apply(files, *packageIdentifier, *rewrite, stdout, stderr)
}

// apply reports or rewrites each file whose header differs.
func apply(files []string, packageIdentifier string, rewrite bool, stdout, stderr io.Writer) int {
	status := 0

	for _, path := range files {
		src, err := os.ReadFile(path) //nolint:gosec // The caller names the path, and reading it is the point.
		if err != nil {
			writef(stderr, "%v\n", err)

			status = 1

			continue
		}

		got, changed := header.Apply(src, packageIdentifier)
		if !changed {
			continue
		}

		if !rewrite {
			writef(stdout, "%s\n", path)

			status = 1

			continue
		}

		err = os.WriteFile(path, got, fileMode) //nolint:gosec // The caller names the path, and writing it is the point.
		if err != nil {
			writef(stderr, "%v\n", err)

			status = 1

			continue
		}

		writef(stdout, "%s\n", path)
	}

	return status
}

// collect returns every Go file under the given paths. A path that names a file
// is returned as it is, and a path that names a directory is walked.
func collect(paths []string) ([]string, error) {
	var files []string

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			files = append(files, path)

			continue
		}

		err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				if skipped[d.Name()] {
					return fs.SkipDir
				}

				return nil
			}

			if strings.HasSuffix(p, ".go") {
				files = append(files, p)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}
