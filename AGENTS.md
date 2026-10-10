# AGENTS.md

Test: `just test` | Before submitting: `just ready`

Read [CONTRIBUTING.md](CONTRIBUTING.md) first. It covers prerequisites, layout,
conventions, and testing for both people and agents.

## Working rules

- Use `fd`, not `find`; `jaq`, not `jq`; and `rg`, not `grep`.
- Never commit unless explicitly asked.
- Never commit directly to main. Work on a branch.
- Use short, descriptive branch names such as `repo-setup`, without personal
  prefixes.
- Repository setup and documentation changes do not authorize moving Go packages
  or changing playbook behavior.
- Keep documentation focused on Voidspan's current implementation.

## Running tools

Use the configured toolchain:

```bash
mise exec -- just fetch
mise exec -- just test
```

From a checkout, run `mise exec -- just run --help`. An installed binary may be
older than the source under development.

Run checks for the change being made. State which passed and which could not
run. For repository setup or documentation changes, use formatting checks and
format only documentation; do not run `just ready` if Go source must stay
unchanged.
