# Development Guide / Руководство разработчика

## Русский

Тесты:

```bash
go test ./cmd/... ./internal/...
cd web && npm test -- --run
```

Линт:

```bash
golangci-lint run ./...
cd web && npm run lint
```

Принципы:
- небольшие и проверяемые PR
- обновление документации вместе с изменениями
- учитывать безопасность на уровне input/auth/config

## English

Tests:

```bash
go test ./cmd/... ./internal/...
cd web && npm test -- --run
```

Lint:

```bash
golangci-lint run ./...
cd web && npm run lint
```

Principles:
- small, testable PRs
- update docs with behavior changes
- treat input/auth/config as security-sensitive by default
