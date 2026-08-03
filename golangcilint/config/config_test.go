/*
 * This file is part of the Valkyrja golangci-lint package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Tests live in an external `_test` package and are co-located with the source
// they cover — the Go convention.
package config_test

import (
	"slices"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/valkyrjaio/ci-golangcilint-go/v26/golangcilint/config"
)

func TestGetReturnsParsableYaml(t *testing.T) {
	t.Parallel()

	var parsed map[string]any
	err := yaml.Unmarshal([]byte(config.Get()), &parsed)
	if err != nil {
		t.Fatalf("the shipped configuration is not valid YAML: %v", err)
	}

	if parsed["version"] != "2" {
		t.Errorf("version = %v, want \"2\"", parsed["version"])
	}
}

func TestGetEndsWithANewline(t *testing.T) {
	t.Parallel()

	if got := config.Get(); !strings.HasSuffix(got, "\n") {
		t.Error("the shipped configuration does not end with a newline")
	}
}

// `goheader` stays disabled. Its check duplicates the copyright header check that
// every repository already runs, and its fix writes a header that the check
// rejects. `valkyrjalint header` writes the header instead.
func TestGetDisablesGoheader(t *testing.T) {
	t.Parallel()

	var parsed struct {
		Linters struct {
			Disable []string `yaml:"disable"`
		} `yaml:"linters"`
	}

	err := yaml.Unmarshal([]byte(config.Get()), &parsed)
	if err != nil {
		t.Fatalf("the shipped configuration is not valid YAML: %v", err)
	}

	if !slices.Contains(parsed.Linters.Disable, "goheader") {
		t.Error("goheader is not in the disable list")
	}
}
