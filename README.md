<p align="center">
  <picture>
    <source srcset="asset/logo-dark.svg" media="(prefers-color-scheme: dark)">
    <source srcset="asset/logo-light.svg" media="(prefers-color-scheme: light)">
    <img src="asset/logo-dark.svg" alt="voidspan" width="480">
  </picture>
</p>

<p align="center">Inspect Ansible-style playbooks with Go.</p>

<p align="center">
  <a href="https://github.com/retr0h/voidspan/releases/latest"><img alt="release" src="https://img.shields.io/github/release/retr0h/voidspan.svg?style=for-the-badge"></a>
  <a href="https://codecov.io/gh/retr0h/voidspan"><img alt="codecov" src="https://img.shields.io/codecov/c/github/retr0h/voidspan?style=for-the-badge"></a>
  <a href="LICENSE"><img alt="license" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge"></a>
  <a href="https://github.com/retr0h/voidspan/actions/workflows/go.yml"><img alt="build" src="https://img.shields.io/github/actions/workflow/status/retr0h/voidspan/go.yml?style=for-the-badge"></a>
  <a href="https://goreleaser.com"><img alt="powered by" src="https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge"></a>
  <a href="https://conventionalcommits.org"><img alt="conventional commits" src="https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge"></a>
  <a href="https://just.systems"><img alt="built with just" src="https://img.shields.io/badge/Built_with-Just-black?style=for-the-badge&logo=just&logoColor=white"></a>
  <img alt="github commit activity" src="https://img.shields.io/github/commit-activity/m/retr0h/voidspan?style=for-the-badge">
  <a href="https://pkg.go.dev/github.com/retr0h/voidspan"><img alt="go reference" src="https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge"></a>
</p>

Voidspan parses [Ansible](https://github.com/ansible/ansible)-style YAML
playbooks and prints their plays, tasks, module arguments, and variables. Task
includes and roles expand into the displayed task list.

The current implementation is a parser prototype. The `run` command does not
execute modules, connect to hosts, or evaluate templates and loops.

## Features

| Feature             | What it does                                                       |
| ------------------- | ------------------------------------------------------------------ |
| Playbook inspection | Displays play names, hosts, tasks, arguments, and task variables   |
| Task includes       | Expands `include_tasks` paths relative to the containing task file |
| Role loading        | Loads `tasks/main.yml` for roles included by a playbook task       |
| Template helper     | Provides an internal Jinja2 renderer, separate from the CLI        |
| Portable builds     | GoReleaser builds Linux and macOS binaries for amd64 and arm64     |

## Quickstart

From a checkout:

```bash
git clone https://github.com/retr0h/voidspan.git
cd voidspan
mise install
just fetch
just deps
just run --help
just build
./dist/voidspan run --playbook testdata/playbook.yml --roles-path testdata/roles
```

Both path flags are required. The checked-in example expands a role and its
included task file, then prints the resulting tasks. The
[usage guide](docs/usage.md) explains the command and its current limits.

Help, usage, and task headings use the same cyan palette as the logo. Use
`--color always` or `--color never` to override terminal detection. `NO_COLOR`
disables automatic color.

## Documentation

- [Usage](docs/usage.md) and [playbook support](docs/playbooks.md)
- [Current architecture](docs/architecture.md)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for prerequisites, setup, conventions,
and testing. The [AI usage policy](AI_POLICY.md) applies to outside
contributions.

## License

The MIT License, see [LICENSE](LICENSE).
