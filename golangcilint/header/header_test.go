/*
 * This file is part of the Valkyrja golangci-lint package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Tests live in an external `_test` package and are co-located with the source
// they cover — the Go convention.
package header_test

import (
	"strings"
	"testing"

	"github.com/valkyrjaio/ci-golangcilint-go/v26/golangcilint/header"
)

const identifier = "Valkyrja golangci-lint"

// want is the exact header the tool writes for `identifier`.
const want = `/*
 * This file is part of the Valkyrja golangci-lint package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */
`

func TestRenderWritesTheSpecifiedHeader(t *testing.T) {
	t.Parallel()

	if got := header.Render(identifier); got != want {
		t.Errorf("Render(%q) =\n%q\nwant\n%q", identifier, got, want)
	}
}

func TestRenderNamesEachPackage(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"Valkyrja Framework", "Sindri", "Project Template"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := header.Render(name)
			if !strings.Contains(got, "This file is part of the "+name+" package.") {
				t.Errorf("Render(%q) does not name the package:\n%s", name, got)
			}
		})
	}
}

func TestApply(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src         string
		want        string
		wantChanged bool
	}{
		"adds the header to a file that carries none": {
			src:         "package constant\n",
			want:        want + "\npackage constant\n",
			wantChanged: true,
		},
		"keeps a package doc comment that sits where the header belongs": {
			src:         "// Package constant holds constants.\npackage constant\n",
			want:        want + "\n// Package constant holds constants.\npackage constant\n",
			wantChanged: true,
		},
		"keeps a block doc comment that is not a header": {
			src:         "/*\n * Notes about this file.\n */\npackage constant\n",
			want:        want + "\n/*\n * Notes about this file.\n */\npackage constant\n",
			wantChanged: true,
		},
		"reports no change where the header is already correct": {
			src:         want + "\npackage constant\n",
			want:        want + "\npackage constant\n",
			wantChanged: false,
		},
		"replaces a header that names the wrong package": {
			src: "/*\n * This file is part of the Project Template package.\n */\n" +
				"\npackage constant\n",
			want:        want + "\npackage constant\n",
			wantChanged: true,
		},
		"replaces a header whose indentation is wrong": {
			src: "/*\n * This file is part of the Valkyrja golangci-lint package.\n" +
				"  *\n  * Copyright (c) 2016-present Melech Mizrachi\n*/\npackage constant\n",
			want:        want + "package constant\n",
			wantChanged: true,
		},
		"adds the header to an empty file": {
			src:         "",
			want:        want + "\n",
			wantChanged: true,
		},
		"adds the header where a block comment never closes": {
			src:         "/*\n * This file is part of the Project Template package.\npackage constant\n",
			want:        want + "\n/*\n * This file is part of the Project Template package.\npackage constant\n",
			wantChanged: true,
		},
		"keeps a build constraint below the header": {
			src:         "//go:build linux\n\npackage constant\n",
			want:        want + "\n//go:build linux\n\npackage constant\n",
			wantChanged: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, changed := header.Apply([]byte(test.src), identifier)

			if string(got) != test.want {
				t.Errorf("Apply() =\n%q\nwant\n%q", got, test.want)
			}

			if changed != test.wantChanged {
				t.Errorf("Apply() changed = %v, want %v", changed, test.wantChanged)
			}
		})
	}
}

// A second Apply over its own output must change nothing. A fix that is not
// idempotent leaves a file that the check rejects, which is the defect that
// `goheader` carries.
func TestApplyIsIdempotent(t *testing.T) {
	t.Parallel()

	sources := []string{
		"package constant\n",
		"// Package constant holds constants.\npackage constant\n",
		"/*\n * This file is part of the Project Template package.\n */\npackage constant\n",
		"",
	}

	for _, src := range sources {
		t.Run(src, func(t *testing.T) {
			t.Parallel()

			once, _ := header.Apply([]byte(src), identifier)

			twice, changed := header.Apply(once, identifier)
			if changed {
				t.Error("Apply() reported a change on its own output")
			}

			if string(twice) != string(once) {
				t.Errorf("Apply() is not idempotent:\n%q\nbecame\n%q", once, twice)
			}
		})
	}
}

// Apply must never drop content. Every byte of the input that is not a header
// survives the call.
func TestApplyKeepsTheDeclaration(t *testing.T) {
	t.Parallel()

	src := "// Package constant holds constants.\npackage constant\n\nconst A = 1\n"

	got, _ := header.Apply([]byte(src), identifier)

	for _, keep := range []string{"// Package constant holds constants.", "package constant", "const A = 1"} {
		if !strings.Contains(string(got), keep) {
			t.Errorf("Apply() dropped %q:\n%s", keep, got)
		}
	}
}
