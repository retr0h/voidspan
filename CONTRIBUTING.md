# Contributing

Thanks for contributing to Voidspan.

## Before you start

Read the [Code of Conduct](CODE_OF_CONDUCT.md) and
[AI usage policy](AI_POLICY.md). Check existing issues and pull requests before
starting a change. Discuss changes to playbook parsing before breaking an
existing configuration.

## Prerequisites

Install [mise](https://mise.jdx.dev/), then install the tools declared in
[.mise.toml](.mise.toml):

```bash
mise install
```

The development toolchain pins Go 1.27.1 and just 1.43.1, plus uv, fd, ripgrep,
and jaq. CI installs the same tools through mise. The application module
declares Go 1.24 as its minimum version.

[just](https://just.systems/) runs development tasks.
[uv](https://docs.astral.sh/uv/) runs the Markdown formatter. Use `fd`, `rg`,
and `jaq` for file searches, text searches, and JSON respectively.

GoReleaser v2 is needed only for `just snapshot`. The release workflow installs
it automatically.

## Setup

```bash
just fetch
just deps
just run --help
```

`just fetch` downloads the pinned shared Markdown recipes from
[osapi-justfiles](https://github.com/osapi-io/osapi-justfiles). Go recipes live
in this repository. Formatter, linter, and coverage tool versions are pinned in
the justfile and run with `go run module@version`, independently of the
application's existing tool dependencies in `go.mod`.

The Jinja2 test suite uses Kluctl's embedded Python distribution. It extracts
that runtime to a temporary directory; a separate Ansible installation is not
required.

## Project structure

```text
main.go              entry point
cmd/                 Cobra commands and help adapter
internal/cli/        Lip Gloss themes, wordmark, and help rendering
internal/ansible/    playbook loading, task expansion, Jinja2 helper, and tests
testdata/            example playbook and role fixtures
asset/               light and dark SVG logos and a PNG export
docs/                Markdown usage and architecture documentation
.github/workflows/   CI and release workflows
.github/repos.json   desired GitHub settings for repo-sync
justfile             development commands
.mise.toml           development tool versions
```

## Code style

Use small packages with clear ownership. An interface belongs with its consumer.
Keep command parsing and presentation separate from playbook logic as new work
is added. Do not move existing packages as part of a documentation change.

Use early returns, wrap errors with `%w`, and separate standard library,
third-party, and local imports. New handwritten Go files carry the MIT notice
from [LICENSE](LICENSE). Preserve generated-file markers on generated code.

```bash
just go-fmt-check
just go-fmt
just go-vet
```

The enabled linters and any existing-code exceptions live in
[.golangci.yml](.golangci.yml). Keep exceptions narrow and remove them when the
corresponding declarations change.

## Documentation

Documentation is ordinary Markdown. There is no site build or Node dependency.

```bash
just md-fmt-check
just md-fmt
```

Write direct prose and describe what exists in Voidspan. Keep each fact in one
document and link to it. The Code of Conduct keeps the upstream wording and
formatting. The formatter excludes `testdata/` to preserve fixture content.

## Testing

```bash
just test
just build
just go-unit
just go-unit-cov
just go_packages=./internal/ansible/... go-unit
```

`just test` checks module consistency and formatting, runs the linter, and runs
tests with the race detector and coverage. Reports go into `.coverage/`.
`.coverignore` excludes command wiring; Codecov uses that same filtered profile.
Coverage targets track the existing baseline.

Tests live beside the code they cover. New public API tests use
`*_public_test.go` in the external test package. Existing Testify suites remain
supported. Use temporary directories for tests that write files.

Before submitting a change:

```bash
just ready
just test
just build
```

`just ready` formats files and runs generators. Review its diff before staging.
For documentation-only work, run `just md-fmt` and the relevant checks without
rewriting Go source.

## Pull requests

Use a branch and Conventional Commit messages such as
`docs: clarify playbook inspection` or `fix: resolve nested task includes`.
Describe the behavior changed, link related issues, and say which checks ran.
Disclose AI assistance as the [policy](AI_POLICY.md) requires.

The workflows check Go, Markdown, justfile formatting, dependencies, and commit
messages. Tags matching `v*` publish release archives through GoReleaser.
`.github/repos.json` records the desired GitHub settings; editing it does not
apply those settings to GitHub by itself.
