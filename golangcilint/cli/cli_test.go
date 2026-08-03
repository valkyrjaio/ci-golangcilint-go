/*
 * This file is part of the Valkyrja golangci-lint package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Tests live in an external `_test` package and are co-located with the source
// they cover — the Go convention.
package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/valkyrjaio/ci-golangcilint-go/v26/golangcilint/cli"
)

const identifier = "Valkyrja golangci-lint"

// headed is a file that already carries the correct header.
const headed = `/*
 * This file is part of the Valkyrja golangci-lint package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant
`

// run calls Run and returns the status with both streams.
func run(args ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer

	status := cli.Run(args, &stdout, &stderr)

	return status, stdout.String(), stderr.String()
}

func TestRunReportsMisuse(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args []string
		want string
	}{
		"no command":          {args: nil, want: "Usage:"},
		"unknown command":     {args: []string{"nope"}, want: `unknown command "nope"`},
		"config takes no arg": {args: []string{"config", "extra"}, want: "config takes no argument"},
		"header needs a name": {args: []string{"header"}, want: "header needs -package"},
		"header bad flag":     {args: []string{"header", "-nope"}, want: "flag provided but not defined"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			status, _, stderr := run(test.args...)

			if status != 2 {
				t.Errorf("status = %d, want 2", status)
			}

			if !strings.Contains(stderr, test.want) {
				t.Errorf("stderr = %q, want it to contain %q", stderr, test.want)
			}
		})
	}
}

func TestConfigWritesTheConfiguration(t *testing.T) {
	t.Parallel()

	status, stdout, stderr := run("config")

	if status != 0 {
		t.Errorf("status = %d, want 0, stderr %q", status, stderr)
	}

	if !strings.Contains(stdout, "version:") {
		t.Errorf("stdout does not hold the configuration:\n%s", stdout)
	}
}

func TestHeaderReportsAFileAndLeavesItAlone(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	write(t, path, "package constant\n")

	status, stdout, _ := run("header", "-package", identifier, dir)

	if status != 1 {
		t.Errorf("status = %d, want 1", status)
	}

	if !strings.Contains(stdout, path) {
		t.Errorf("stdout = %q, want it to name %q", stdout, path)
	}

	if got := read(t, path); got != "package constant\n" {
		t.Errorf("the file changed without -w:\n%q", got)
	}
}

func TestHeaderWritesTheHeader(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	write(t, path, "package constant\n")

	status, _, stderr := run("header", "-package", identifier, "-w", dir)

	if status != 0 {
		t.Errorf("status = %d, want 0, stderr %q", status, stderr)
	}

	if got := read(t, path); got != headed {
		t.Errorf("the file reads\n%q\nwant\n%q", got, headed)
	}
}

func TestHeaderIsQuietWhereEveryHeaderIsCorrect(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.go"), headed)

	status, stdout, _ := run("header", "-package", identifier, dir)

	if status != 0 {
		t.Errorf("status = %d, want 0", status)
	}

	if stdout != "" {
		t.Errorf("stdout = %q, want it empty", stdout)
	}
}

func TestHeaderReadsAFilePathAndSkipsANonGoFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	goFile := filepath.Join(dir, "a.go")
	write(t, goFile, headed)
	write(t, filepath.Join(dir, "notes.txt"), "no header here\n")

	status, stdout, _ := run("header", "-package", identifier, goFile, dir)

	if status != 0 {
		t.Errorf("status = %d, want 0, stdout %q", status, stdout)
	}
}

func TestHeaderSkipsAnExcludedDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	vendor := filepath.Join(dir, "vendor")

	err := os.Mkdir(vendor, 0o750)
	if err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	write(t, filepath.Join(vendor, "dep.go"), "package dep\n")

	status, stdout, _ := run("header", "-package", identifier, dir)

	if status != 0 {
		t.Errorf("status = %d, want 0; it read a vendored file: %q", status, stdout)
	}
}

func TestHeaderReportsAMissingPath(t *testing.T) {
	t.Parallel()

	status, _, stderr := run("header", "-package", identifier, filepath.Join(t.TempDir(), "absent"))

	if status != 1 {
		t.Errorf("status = %d, want 1", status)
	}

	if stderr == "" {
		t.Error("stderr is empty, want the stat failure")
	}
}

func TestHeaderDefaultsToTheWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.go"), "package constant\n")

	t.Chdir(dir)

	status, stdout, _ := run("header", "-package", identifier)

	if status != 1 {
		t.Errorf("status = %d, want 1", status)
	}

	if !strings.Contains(stdout, "a.go") {
		t.Errorf("stdout = %q, want it to name a.go", stdout)
	}
}

func TestHeaderReportsAFileItCannotRead(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	write(t, path, "package constant\n")

	err := os.Chmod(path, 0o000)
	if err != nil {
		t.Fatalf("chmod: %v", err)
	}

	t.Cleanup(func() { _ = os.Chmod(path, 0o600) })

	status, _, stderr := run("header", "-package", identifier, dir)

	if status != 1 {
		t.Errorf("status = %d, want 1", status)
	}

	if stderr == "" {
		t.Error("stderr is empty, want the read failure")
	}
}

func TestHeaderReportsAFileItCannotWrite(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	write(t, path, "package constant\n")

	err := os.Chmod(path, 0o400)
	if err != nil {
		t.Fatalf("chmod: %v", err)
	}

	t.Cleanup(func() { _ = os.Chmod(path, 0o600) })

	status, _, stderr := run("header", "-package", identifier, "-w", dir)

	if status != 1 {
		t.Errorf("status = %d, want 1", status)
	}

	if stderr == "" {
		t.Error("stderr is empty, want the write failure")
	}
}

func TestHeaderReportsADirectoryItCannotWalk(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	closed := filepath.Join(dir, "closed")

	err := os.Mkdir(closed, 0o750)
	if err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	write(t, filepath.Join(closed, "a.go"), "package constant\n")

	err = os.Chmod(closed, 0o000)
	if err != nil {
		t.Fatalf("chmod: %v", err)
	}

	t.Cleanup(func() { _ = os.Chmod(closed, 0o750) })

	status, _, stderr := run("header", "-package", identifier, dir)

	if status != 1 {
		t.Errorf("status = %d, want 1", status)
	}

	if stderr == "" {
		t.Error("stderr is empty, want the walk failure")
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()

	err := os.WriteFile(path, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func read(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return string(content)
}
