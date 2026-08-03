#
# This file is part of the Project Template package.
#
# Copyright (c) 2016-present Melech Mizrachi
#
# Released under the MIT License. See LICENSE.md for details.
#

# Root task runner for the Valkyrja Go template — the analog of PHP's
# composer.json, TypeScript's package.json, and Java's Gradle scripts. Drive
# every CI tool through these targets; check here first for exact target names.
#
# golangci-lint is pinned in .github/ci/lint/go.mod (an isolated tool module) and
# run via `go tool` so its version never mixes with the framework's own deps.
GOLANGCI_LINT ?= go tool -modfile=.github/ci/lint/go.mod golangci-lint

# The coverage floor, as a percentage. 100 is the definition of done; it is a hard
# floor, never lowered to accommodate a gap (cover the code, or leave it out of the
# build). Override only for a local spot check: `make coverage COVERAGE_FLOOR=90`.
COVERAGE_FLOOR ?= 100

.DEFAULT_GOAL := ci

.PHONY: ci
ci: tidy-check fmt-check lint coverage ## Run the full CI gate

.PHONY: build
build: ## Compile every package
	go build ./...

.PHONY: test
test: ## Run tests with the race detector
	go test -race ./...

.PHONY: coverage
coverage: ## Run tests with coverage and fail below COVERAGE_FLOOR
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
# `go tool cover` only reports and always exits 0, so without this the profile was
# generated and then ignored — a run at 55% passed exactly like one at 100%.
#
# Statement coverage only: Go has no native branch coverage, so an untested branch
# inside a covered statement stays invisible here and is a review concern (see
# architecture/go/AGENTS.md).
#
# A package of nothing but constants has no statements at all, and `go tool cover`
# prints "total: (statements) 0.0%" for that — identical to genuinely covering
# nothing. Counting the per-function rows tells the two apart: no rows means there
# was nothing to cover, which is a pass, not a 0% failure.
	@go tool cover -func=coverage.out | awk -v floor='$(COVERAGE_FLOOR)' ' \
		/^total:/ { total = $$3; next } \
		{ rows++ } \
		END { \
			if (rows == 0) { print "No statements to cover; nothing to enforce."; exit 0 } \
			sub(/%$$/, "", total); \
			if (total + 0 < floor + 0) { \
				printf "FAIL  statement coverage %.1f%% < %.1f%%\n", total, floor; exit 1 \
			} \
			printf "PASS  statement coverage %.1f%% >= %.1f%%\n", total, floor \
		}'

.PHONY: fmt
fmt: ## Format the codebase
	$(GOLANGCI_LINT) fmt

.PHONY: fmt-check
fmt-check: ## Fail if the codebase is not formatted
	$(GOLANGCI_LINT) fmt --diff

.PHONY: lint
lint: ## Run golangci-lint (static analysis, security, architecture, dead code)
	$(GOLANGCI_LINT) run

.PHONY: tidy
tidy: ## Tidy module dependencies
	go mod tidy

.PHONY: tidy-check
tidy-check: ## Fail if go.mod / go.sum are not tidy
	go mod tidy -diff

.PHONY: verify
verify: ## Verify dependencies against go.sum
	go mod verify

.PHONY: tools
tools: ## Download the pinned CI tool dependencies
	cd .github/ci/lint && go mod download

.PHONY: outdated
outdated: ## List direct dependencies with newer versions available
	go list -u -m -f '{{if and .Update (not .Indirect)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' all

.PHONY: update
update: ## Update all dependencies to latest and tidy
	go get -u ./...
	go mod tidy

.PHONY: help
help: ## List available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
