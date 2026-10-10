set allow-duplicate-variables := true

import? '.just/remote/md.just'

md_site_dir := ""
md_extra_excludes := "--exclude CODE_OF_CONDUCT.md --exclude 'testdata/**'"
justfiles_ref := "7bf2cba5c9f39e3d3ffca835fa2c85ed42a67f09"

# Match Psion's development tools without changing application dependencies.

gofumpt := "go run mvdan.cc/gofumpt@v0.12.0"
golines := "go run github.com/segmentio/golines@v0.13.0"
golangci := "go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0"
go_packages := "./..."

# List development commands
default:
    @just --list

# Fetch the shared Markdown recipes
fetch:
    mkdir -p .just/remote
    curl -sSfL https://raw.githubusercontent.com/osapi-io/osapi-justfiles/{{ justfiles_ref }}/md/md.just -o .just/remote/md.just

# Download dependencies and check that development tools resolve
deps:
    go mod download
    {{ gofumpt }} -version
    {{ golines }} --version
    {{ golangci }} version

# Compile and run the checkout
run *args:
    go run main.go {{ args }}

# Build the current platform's binary
build:
    mkdir -p dist
    go build -o dist/voidspan .

# Build release archives locally without publishing
snapshot:
    goreleaser release --snapshot --clean

# Download and tidy application dependencies when intentionally maintaining them
go-mod:
    go mod download
    go mod tidy

# Check module consistency without modifying go.mod or go.sum
go-mod-check:
    go mod tidy -diff

# Format handwritten Go files
go-fmt:
    fd -t f -e go -E '*.gen.go' -E '*.pb.go' -0 | xargs -0 {{ gofumpt }} -w
    fd -t f -e go -E '*.gen.go' -E '*.pb.go' -0 | xargs -0 {{ golines }} --base-formatter='{{ gofumpt }}' -w

# Check Go formatting without rewriting source
go-fmt-check:
    #!/usr/bin/env bash
    set -euo pipefail
    changes=$(fd -t f -e go -E '*.gen.go' -E '*.pb.go' -0 | xargs -0 {{ gofumpt }} -l)
    if [ -n "$changes" ]; then echo "$changes"; exit 1; fi
    changes=$(fd -t f -e go -E '*.gen.go' -E '*.pb.go' -0 | xargs -0 {{ golines }} --dry-run --base-formatter='{{ gofumpt }}' -l)
    if [ -n "$changes" ]; then echo "$changes"; exit 1; fi

# Run the repository's configured linters
go-vet:
    {{ golangci }} run --config .golangci.yml

# Run the unit tests with the race detector
go-unit:
    go test -race {{ go_packages }}

# Write the coverage profile, excluding only command wiring
go-unit-cov:
    mkdir -p .coverage
    go test -race -coverprofile=.coverage/cover.raw.out {{ go_packages }}
    rg -v -f .coverignore .coverage/cover.raw.out > .coverage/cover.out
    go run github.com/boumenot/gocover-cobertura@v1.5.0 < .coverage/cover.out > .coverage/cobertura.xml
    go tool cover -func=.coverage/cover.out

# Open the coverage report
go-unit-cov-map: go-unit-cov
    go tool cover -html=.coverage/cover.out

# Run source checks and tests
test: go-mod-check go-fmt-check go-vet go-unit-cov

# Run any Go generators
generate:
    go generate ./...

# Format the justfile
just-fmt:
    just --fmt --unstable

# Check the justfile's formatting
just-fmt-check:
    just --fmt --check --unstable

# Format and lint before submitting a change
ready:
    just generate
    just md-fmt
    just go-fmt
    just go-vet
    just just-fmt
