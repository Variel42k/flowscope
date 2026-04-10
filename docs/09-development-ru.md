# Руководство разработчика (RU)

## Локальная разработка

### Backend

Основные сервисы:

- `cmd/collector`
- `cmd/api`
- `cmd/worker`
- `cmd/seed`
- `cmd/flowgen`

Ключевые пакеты:

- `internal/decoder`
- `internal/normalize`
- `internal/enrich`
- `internal/storage`
- `internal/graph`
- `internal/api`

### Frontend

Папка: `web`

Стек:

- React + TypeScript + Vite
- Tailwind
- ECharts
- Cytoscape.js
- TanStack Table

## Тесты

```bash
go test ./cmd/... ./internal/...
cd web && npm test -- --run
```

## Линтинг/форматирование

```bash
golangci-lint run ./...
cd web && npm run lint
cd web && npm run format
```

## Принципы изменений

1. Сначала small working change.
2. Потом тесты и документация.
3. Отдельные PR для:
   - функциональности
   - миграций схемы
   - UI-изменений
   - security hardening

## Изменения схемы ClickHouse

- SQL файлы лежат в `deploy/clickhouse/init`.
- Для обратной совместимости придерживайтесь add-first миграций.

## Рекомендации по производительности

1. Использовать rollup tables для дашбордов/графа.
2. Использовать raw table только для drill-down.
3. Ограничивать time windows в тяжёлых запросах.
