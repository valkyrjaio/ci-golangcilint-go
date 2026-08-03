/*
 * This file is part of the Valkyrja golangci-lint package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package header renders the Valkyrja copyright header and applies it to a Go
// source file.
//
// COPYRIGHT_HEADER.md in the .github repository specifies the text. This package
// holds the only Go copy of it, so a repository that consumes this module states
// its package identifier and nothing else.
//
// golangci-lint's `goheader` linter checks the header and cannot add one. Its
// fix writes `// ` line comments whenever a file carries no header, because
// `generateFix` prefixes every line with `// ` in that branch. It also replaces
// the first comment in the file, which deletes a package doc comment that sits
// where the header belongs. This package writes the block form, and it replaces
// a comment only where that comment is a copyright header.
package header

import (
	"bytes"
	"strings"
)

// marker identifies a copyright header. A comment that holds this text is a
// header the tool may replace. Every other comment is documentation, and the
// tool leaves it alone.
const marker = "This file is part of the "

// Render returns the copyright header for a package identifier. The result ends
// with a newline, and it carries no blank line after it.
func Render(packageIdentifier string) string {
	var b strings.Builder

	b.WriteString("/*\n")
	b.WriteString(" * " + marker + packageIdentifier + " package.\n")
	b.WriteString(" *\n")
	b.WriteString(" * Copyright (c) 2016-present Melech Mizrachi\n")
	b.WriteString(" *\n")
	b.WriteString(" * Released under the MIT License. See LICENSE.md for details.\n")
	b.WriteString(" */\n")

	return b.String()
}

// Apply returns src with the header for a package identifier, and reports
// whether the content changed.
//
// Apply adds the header where the file carries none. It replaces the header
// where the file opens with a copyright header that differs. It leaves every
// other comment in place, so a file that opens with a package doc comment keeps
// that comment.
func Apply(src []byte, packageIdentifier string) ([]byte, bool) {
	want := []byte(Render(packageIdentifier))

	if bytes.HasPrefix(src, want) {
		return src, false
	}

	if end, ok := headerBlockEnd(src); ok {
		return append(append([]byte{}, want...), src[end:]...), true
	}

	return append(append(append([]byte{}, want...), '\n'), src...), true
}

// headerBlockEnd reports the offset just past the opening block comment where
// that comment is a copyright header. A file that opens with any other comment
// has no such block, so the caller adds the header instead of replacing one.
func headerBlockEnd(src []byte) (int, bool) {
	if !bytes.HasPrefix(src, []byte("/*")) {
		return 0, false
	}

	closeIndex := bytes.Index(src, []byte("*/"))
	if closeIndex < 0 {
		return 0, false
	}

	end := closeIndex + len("*/")

	if !bytes.Contains(src[:end], []byte(marker)) {
		return 0, false
	}

	// Consume the newline that ends the comment line, so the replacement does not
	// leave a blank line where the old header sat.
	if end < len(src) && src[end] == '\n' {
		end++
	}

	return end, true
}
