<p align="center"><a href="https://valkyrja.io" target="_blank">
    <img src="https://raw.githubusercontent.com/valkyrjaio/art/refs/heads/master/long-banner/orange/go.png" width="100%">
</a></p>

# Project Template (Go)

A starter template for creating new Go repositories in the Valkyrjaio
organization.

This template ships with the full Valkyrja CI pipeline pre-wired (golangci-lint
and the built-in Go test/coverage toolchain), a minimal Go module setup, and the
repository conventions used across the rest of the org. Use it as the starting
point for any new Go package, CI tool config, or integration repo — not for
end-user applications built on the Valkyrja framework (use
[`valkyrja-starter-app-go`][starter url] for that).

<p>
    <a href="https://pkg.go.dev/github.com/valkyrjaio/project-template-go/v26"><img src="https://pkg.go.dev/badge/github.com/valkyrjaio/project-template-go/v26.svg" alt="Go Reference"></a>
    <a href="https://github.com/valkyrjaio/project-template-go/blob/26.x/LICENSE.md"><img src="https://img.shields.io/github/license/valkyrjaio/project-template-go.svg" alt="License"></a>
    <a href="https://github.com/valkyrjaio/project-template-go/actions/workflows/ci.yml?query=branch%3A26.x"><img src="https://github.com/valkyrjaio/project-template-go/actions/workflows/ci.yml/badge.svg?branch=26.x" alt="CI Status"></a>
</p>

Usage
-----

### Use this template _(recommended)_

This repository is a GitHub template. Click the **Use this template** button
at the top of the repo to create a new repository in the Valkyrjaio
organization, pre-populated with the template's structure and CI.

### After Creating Your Repo

1. Update `go.mod` with your module path (keep the `/vNN` major-version suffix)
2. Rename the `template/` package directory to your package and replace its
   contents with your source code
3. Update this `README.md` to describe the new package
4. Configure the required secrets and variables — see
   [`REPOSITORY_NAMING.md`][repository naming url] for naming guidance and
   `.github`'s workflow documentation for secret requirements
5. Verify CI passes on the first commit

What's Included
---------------

- **CI pipeline** — golangci-lint (formatting, static analysis, security,
  architecture, dead code) plus the built-in `go test` toolchain for tests and
  coverage, matching the org convention
- **Makefile** — the root task runner (`make ci`, `make lint`, `make test`,
  `make coverage`, `make fmt`), the Go analog of composer/npm/Gradle scripts
- **Repository conventions** — aligned with
  [`REPOSITORY_NAMING.md`][repository naming url] and
  [`VOCABULARY.md`][vocabulary url]

### CI tooling

Go's toolchain covers most of what other languages need separate tools for
(testing, coverage, formatting, vet, build, dependency hygiene). The only
external tool is **golangci-lint**, pinned in `.github/ci/lint/go.mod` and run
via `make lint`. Run the full gate locally with `make ci`.

Versioning and Release Process
------------------------------

This template follows [semantic versioning][semantic versioning url] with a
major release every year, and support for each major version for 2 years
from the date of release.

Go [semantic import versioning][semantic import versioning url] requires the
major version to be encoded in the module path for any major `>= 2` — hence the
`/v26` suffix in `go.mod`. Each annual major bump advances that suffix.

For more information see our
[Versioning and Release Process documentation][Versioning and Release Process url].

### Supported Versions

Bug fixes are provided until 3 months after the next major release. Security
fixes are provided for 2 years after the initial release.

| Version | Go    | Release        | Bug Fixes Until | Security Fixes Until |
| :------ | :---- | :------------- | :-------------- | :------------------- |
| 26      | 1.26+ | March 31, 2026 | Q2 2027         | Q1 2028              |

Contributing
------------

This template is an open-source, community-driven project. Improvements to
the template itself — refinements to the included CI configuration, Go module
setup, or documentation — are welcome.

See [`CONTRIBUTING.md`][contributing url] for the submission process and
[`VOCABULARY.md`][vocabulary url] for the terminology used across Valkyrja.

Security Issues
---------------

If you discover a security vulnerability, please follow our
[disclosure procedure][security vulnerabilities url].

License
-------

This template is open-source software licensed under the
[MIT license][MIT license url]. See [`LICENSE.md`](./LICENSE.md).

[Valkyrja url]: https://valkyrja.io
[starter url]: https://github.com/valkyrjaio/valkyrja-starter-app-go
[repository naming url]: https://github.com/valkyrjaio/.github/blob/26.x/REPOSITORY_NAMING.md
[vocabulary url]: https://github.com/valkyrjaio/.github/blob/26.x/VOCABULARY.md
[contributing url]: https://github.com/valkyrjaio/.github/blob/26.x/CONTRIBUTING.md
[security vulnerabilities url]: https://github.com/valkyrjaio/.github/blob/26.x/SECURITY.md
[Versioning and Release Process url]: https://github.com/valkyrjaio/.github/blob/26.x/VERSIONING_AND_RELEASE_PROCESS.md
[semantic versioning url]: https://semver.org/
[semantic import versioning url]: https://go.dev/ref/mod#major-version-suffixes
[MIT license url]: https://opensource.org/licenses/MIT
[license url]: ./LICENSE.md
