// The `/vN` path suffix is required by Go's semantic import versioning: any
// major version >= 2 must be encoded in the module path. It tracks the annual
// major version (VERSION.md) and is bumped each major release — the version
// branch workflow rewrites the suffix here and in every import of it.
module github.com/valkyrjaio/ci-golangcilint-go/v26

go 1.26

toolchain go1.26.0

require gopkg.in/yaml.v3 v3.0.1
