/*
 * This file is part of the Project Template package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Tests live in an external `_test` package and are co-located with the source
// they cover — the Go convention. Reusable doubles belong in a `fixtures`
// package named `*Fixture`, never `*Test`.
//
// The release workflow rewrites both constants. Each test asserts a format and
// never an exact value.
package constant_test

import (
	"regexp"
	"testing"

	"github.com/valkyrjaio/project-template-go/v26/template/constant"
)

// versionPattern is the MAJOR.MINOR.PATCH format that the release workflow writes.
var versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

// versionBuildDateTimePattern is the `Month D YYYY HH:MM:SS MST` format that the
// release workflow writes with `date '+%B %-d %Y %T MST'`.
var versionBuildDateTimePattern = regexp.MustCompile(`^[A-Z][a-z]+ \d{1,2} \d{4} \d{2}:\d{2}:\d{2} MST$`)

func TestVersionHasTheVersionFormat(t *testing.T) {
	t.Parallel()

	if !versionPattern.MatchString(constant.Version) {
		t.Errorf("Version must match %s, but is: %s", versionPattern, constant.Version)
	}
}

func TestVersionBuildDateTimeHasTheBuildDateTimeFormat(t *testing.T) {
	t.Parallel()

	if !versionBuildDateTimePattern.MatchString(constant.VersionBuildDateTime) {
		t.Errorf(
			"VersionBuildDateTime must match %s, but is: %s",
			versionBuildDateTimePattern,
			constant.VersionBuildDateTime,
		)
	}
}
