/*
 * This file is part of the Valkyrja golangci-lint package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package main

// This test sits in the source package, because it replaces the package-private
// `exit`. It is the one test outside an external `_test` package, and it exists
// so `main` carries a test rather than an exclusion from the coverage report.

import (
	"os"
	"testing"
)

// binary is the name the shell uses, which `os.Args[0]` carries.
const binary = "valkyrjalint"

func TestMainExitsWithTheCommandStatus(t *testing.T) {
	tests := map[string]struct {
		args []string
		want int
	}{
		"a known command succeeds":  {args: []string{binary, "config"}, want: 0},
		"no command reports misuse": {args: []string{binary}, want: 2},
		"an unknown command is 2":   {args: []string{binary, "nope"}, want: 2},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			originalExit := exit
			originalArgs := os.Args
			originalStdout := os.Stdout

			t.Cleanup(func() {
				exit = originalExit
				os.Args = originalArgs
				os.Stdout = originalStdout
			})

			devNull, err := os.Open(os.DevNull)
			if err != nil {
				t.Fatalf("open %s: %v", os.DevNull, err)
			}

			t.Cleanup(func() { _ = devNull.Close() })

			os.Stdout = devNull

			got := -1
			exit = func(code int) { got = code }
			os.Args = test.args

			main()

			if got != test.want {
				t.Errorf("main() exited %d, want %d", got, test.want)
			}
		})
	}
}
