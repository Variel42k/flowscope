# Contributing to FlowScope

Thanks for your interest in contributing.

## Development Setup

1. Clone the repository.
2. Start dependencies:

```bash
docker compose up --build -d clickhouse seed
```

3. Run all services:

```bash
docker compose up --build -d
```

4. Run tests:

```bash
make test
```

## Coding Standards

- Go: keep code simple, explicit, and testable.
- Frontend: TypeScript strictness, no `any` unless justified.
- Prefer small focused changes over large mixed refactors.
- Add tests for behavior changes.
- Keep security in mind for all input paths and config defaults.

## Pull Request Rules

1. Use a descriptive PR title.
2. Explain what changed and why.
3. Include test evidence.
4. Update docs if behavior/config changed.
5. Do not include secrets or private data.

## Commit Style

Use concise imperative commit messages, for example:

- `api: add max body size for inventory import`
- `web: fix flow table pagination total`
- `docs: add deployment runbook`

## Reporting Issues

- Use GitHub Issues templates.
- Include steps to reproduce and expected behavior.
- For sensitive security issues, use [SECURITY.md](SECURITY.md).
