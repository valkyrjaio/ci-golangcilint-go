<p align="center"><a href="https://valkyrja.io" target="_blank">
    <img src="https://raw.githubusercontent.com/valkyrjaio/art/refs/heads/26.x/long-banner/orange/go.png" width="100%">
</a></p>

# Valkyrja golangci-lint (Go)

The golangci-lint configuration and the copyright header tool that every
Valkyrja Go repository runs.

This package holds the linting rules in one place, the way
[`ci-phpcsfixer-php`][phpcsfixer url] holds them for PHP. A consuming repository
states its own package identifier and nothing else.

<p>
    <a href="https://pkg.go.dev/github.com/valkyrjaio/ci-golangcilint-go/v26"><img src="https://pkg.go.dev/badge/github.com/valkyrjaio/ci-golangcilint-go/v26.svg" alt="Go Reference"></a>
    <a href="https://github.com/valkyrjaio/ci-golangcilint-go/releases"><img src="https://img.shields.io/github/v/release/valkyrjaio/ci-golangcilint-go" alt="Latest Version"></a>
    <a href="https://github.com/valkyrjaio/ci-golangcilint-go/blob/26.x/LICENSE.md"><img src="https://img.shields.io/github/license/valkyrjaio/ci-golangcilint-go.svg" alt="License"></a>
    <a href="https://github.com/valkyrjaio/ci-golangcilint-go/actions/workflows/ci.yml?query=branch%3A26.x"><img src="https://github.com/valkyrjaio/ci-golangcilint-go/actions/workflows/ci.yml/badge.svg?branch=26.x" alt="CI Status"></a>
    <a href="https://coveralls.io/github/valkyrjaio/ci-golangcilint-go?branch=26.x"><img src="https://coveralls.io/repos/github/valkyrjaio/ci-golangcilint-go/badge.svg?branch=26.x" alt="Coverage Status"></a>
    <a href="https://sonarcloud.io/summary/new_code?id=valkyrjaio_ci-golangcilint-go"><img src="https://sonarcloud.io/api/project_badges/measure?project=valkyrjaio_ci-golangcilint-go&metric=sqale_rating" alt="Maintainability Rating"></a>
</p>

What This Package Ships
-----------------------

| Command                                | What it does                                          |
| -------------------------------------- | ----------------------------------------------------- |
| `valkyrjalint config`                  | Writes the shared `.golangci.yml` to stdout           |
| `valkyrjalint header -package NAME`    | Reports each Go file whose header is missing or wrong |
| `valkyrjalint header -package NAME -w` | Writes the header into each such file                 |

Install
-------

Add the command to the repository's tool module, next to golangci-lint. That
module is isolated, so its dependency graph never mixes with the repository's
own:

```bash
cd .github/ci/lint && go get -tool github.com/valkyrjaio/ci-golangcilint-go/v26/cmd/valkyrjalint
```

Use
---

Add the targets to the repository's `Makefile`. `PACKAGE_IDENTIFIER` is the value
that [`COPYRIGHT_HEADER.md`][header url] maps to the repository:

```make
VALKYRJALINT ?= go tool -modfile=.github/ci/lint/go.mod valkyrjalint
PACKAGE_IDENTIFIER ?= Valkyrja Framework

.PHONY: config-write
config-write:
	$(VALKYRJALINT) config > .golangci.yml

# Warning: never pipe the generator straight into `diff`. `make` runs a recipe under
# `/bin/sh`, which has no `pipefail`, so the pipeline reports `diff`'s status. A
# generator that fails writes nothing, and the check then blames a stale file for a
# crash. Write the output to a file first, so the generator's own status ends the
# recipe.
.PHONY: config-check
config-check:
	@generated=$$(mktemp); trap 'rm -f "$$generated"' EXIT; \
		$(VALKYRJALINT) config > "$$generated" \
		&& diff -u .golangci.yml "$$generated"

.PHONY: header-fix
header-fix:
	$(VALKYRJALINT) header -package '$(PACKAGE_IDENTIFIER)' -w .

.PHONY: header-check
header-check:
	@$(VALKYRJALINT) header -package '$(PACKAGE_IDENTIFIER)' .
```

Commit the generated `.golangci.yml`. An editor reads that file, so a repository
that wrote it only in CI would leave a developer with no diagnostics.
`config-check` keeps the committed file current, the same way `go mod tidy -diff`
keeps `go.mod` current.

Why The Header Tool Exists
--------------------------

golangci-lint ships a `goheader` linter. It checks a header, and it cannot add
one.

Warning: never enable `goheader` to get an automatic header. Its fix damages a
file, and it reports `0 issues` each time:

- **A file with no header gets `// ` line comments.** `generateFix` prefixes
  every template line with `// ` in that branch, so the tool cannot write the
  block form that every Valkyrja language uses.
- **A package doc comment is deleted.** The fix replaces the first comment in
  the file. In a file that carries no header, that comment is the documentation.
- **A wrong header gets wrong indentation.** The fix adds a space to each line
  after the first, and it removes the space before the closing `*/`. The next
  plain run then fails on the same file.

`valkyrjalint header` writes the block form. It replaces a comment only where
that comment is a copyright header, so a package doc comment survives. A second
run reports no change.

This is why the shipped configuration disables `goheader`.

The copyright header check in the [`.github`][github url] repository still reads
every tracked file, whatever the language. This tool writes the header, and that
check verifies it.

Contributing
------------

Read [`AGENTS.md`](AGENTS.md) and the guides it links.

```bash
make ci
```

[phpcsfixer url]: https://github.com/valkyrjaio/ci-phpcsfixer-php
[header url]: https://github.com/valkyrjaio/.github/blob/26.x/COPYRIGHT_HEADER.md
[github url]: https://github.com/valkyrjaio/.github
