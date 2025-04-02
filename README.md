# Voidspan

Voidspan is a Go-based interpreter for [Ansible][] playbooks by translating
[Ansible][]-flavored YAML into executable Go workflows — no Python runtime required.

## Testing

Enable [Remote Taskfile][] feature.

```bash
export TASK_X_REMOTE_TASKFILES=1
```

Install dependencies:

```bash
$ task go:deps
```

To execute tests:

```bash
$ task go:test
```

Auto format code:

```bash
$ task go:fmt
```

List helpful targets:

```bash
$ task
```

## License

The [MIT][] License.

The logo is licensed under the [Creative Commons NoDerivatives 4.0 License][],
and designed by [@nanotron][].
If you have some other use in mind, contact us.

[Ansible]: https://github.com/ansible/ansible
[Remote Taskfile]: https://taskfile.dev/experiments/remote-taskfiles/
[MIT]: LICENSE
