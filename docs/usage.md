# Usage

Voidspan currently loads playbooks and prints the parsed tasks. It does not
execute Ansible modules or change the target hosts.

## Build and run

Follow the [development setup](../CONTRIBUTING.md#setup), then build a binary:

```bash
just build
./dist/voidspan --help
./dist/voidspan run --help
```

Inspect the checked-in example:

```bash
./dist/voidspan run --playbook testdata/playbook.yml --roles-path testdata/roles
```

To run directly from the checkout:

```bash
just run run --playbook testdata/playbook.yml --roles-path testdata/roles
```

The first `run` selects the just recipe; the second selects Voidspan's command.

## Flags

| Flag           | Short form | Meaning                                    |
| -------------- | ---------- | ------------------------------------------ |
| `--playbook`   | `-p`       | Path to the playbook YAML file             |
| `--roles-path` | `-r`       | Base directory containing role directories |

Both flags must be supplied, even though help displays defaults. The positional
`[playbook.yml]` shown in help is currently ignored; use `--playbook`.

The global `--color` flag accepts `auto`, `always`, or `never`. Automatic color
uses the output terminal and respects `NO_COLOR` and `TERM=dumb`. Piped output
is plain text unless `--color always` is supplied.

## Output

The example prints a play named `converge`, followed by tasks from the playbook,
the `retr0h.role-1` role, and its `include-1.yml` file. Each task shows its
module, raw arguments, task variables, and an optional loop expression.

Templates such as `{{ item }}` and `{{ role_1_items }}` remain literal strings.
Variables attached to an `include_role` task are not propagated to the expanded
role tasks. See [playbook support](playbooks.md) for the current boundaries.

## Release archives

The release configuration builds Linux and macOS archives for amd64 and arm64.
With GoReleaser v2 installed, build a local snapshot without publishing:

```bash
just snapshot
```

Artifacts are written to `dist/`. The GitHub release workflow publishes them
when a tag matching `v*` is pushed.
