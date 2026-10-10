# Playbook support

Voidspan implements a subset of Ansible's playbook structure. It does not
provide Ansible execution or full compatibility.

## Plays and tasks

A playbook is a YAML list of plays. The loader reads each play's string `name`
and `hosts` fields and its `tasks` list. Plays without a task list are skipped.

A task records its name, module, raw arguments, task-level `vars`, string `loop`
expression, and source path. Mapping arguments are retained as a map; scalar
arguments are stored under `__value__`.

Only `name`, `vars`, `loop`, and task includes receive special treatment during
task parsing. Other keys are treated as module names, so execution directives
such as `when`, `register`, and `become` are not supported. Literal list loops
are not captured; string loop expressions are retained without evaluation.

## Task includes

Both `include_tasks` and `ansible.builtin.include_tasks` accept a string path.
The loader reads that file relative to the file containing the include and
recursively expands its task list. The include itself is replaced by its tasks.
Template expressions in include paths are not evaluated.

## Roles

At the playbook task level, `include_role` and `ansible.builtin.include_role`
load `<roles-path>/<name>/tasks/main.yml` and insert those tasks in place of the
include. Task includes inside that role are resolved relative to their files.

The loader does not read role defaults, vars, handlers, or metadata, propagate
the include task's variables, expand role includes nested inside a role, or
process the play-level `roles:` list.

## Templates

`RenderJinjaFields` is an internal helper backed by Kluctl's Jinja2 library and
embedded Python runtime. It renders string values in maps and nested maps with
the supplied context. Other values, including lists, pass through unchanged.

The CLI does not call this helper. Its printed task arguments and loops remain
unrendered.
