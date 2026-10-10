# Architecture

Voidspan is a Go CLI with a playbook parser and a separate template helper.

```text
main.go
   |
cmd.Execute (Cobra)
   |
run --playbook ... --roles-path ...
   |
internal/ansible.LoadPlaybook
   |-- parseTasks: expand task includes
   |-- LoadRoleTasks: read role tasks/main.yml
   |
[]Play / []Task
   |
print plays and task fields
```

`cmd/run.go` reads the two path flags through Viper, loads the playbook from
disk, calls the parser, and prints the returned data. There is no module
executor, remote transport, inventory resolver, or persistent state layer.

`internal/ansible/load_playbook.go` decodes YAML plays, parses task lists, and
expands role includes. `parse_tasks.go` resolves task includes relative to the
source file. `load_role_tasks.go` reads a role's `tasks/main.yml` and passes it
through the same task parser. `types.go` defines `Play` and `Task`.

`render_jinja_fields.go` provides `RenderJinjaFields`, using Kluctl's Jinja2
wrapper. It is covered by tests but is not wired into the command. Kluctl
supplies an embedded Python runtime for this helper.

`internal/cli/` owns the Lip Gloss palette, block wordmark, color detection, and
help renderer. `cmd/help.go` adapts Cobra commands to plain help data, and
`cmd/root.go` installs the renderer. The SVG assets use the same wordmark;
`cmd/run.go` uses the shared theme for task headings and labels.

Tests live beside the parser and use temporary files and Testify suites. The
example playbook and role tree under `testdata/` exercise CLI inspection.
