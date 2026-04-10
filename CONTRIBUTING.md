# Contributing / Как вносить вклад

## Русский

Спасибо за интерес к FlowScope.

### Быстрый старт разработки

```bash
docker compose up --build -d clickhouse seed
docker compose up --build -d
make test
```

### Правила PR

1. Понятный заголовок и описание "что и зачем".
2. Приложить результаты тестов.
3. Обновить документацию при изменении поведения.
4. Не добавлять секреты и приватные данные.

## English

Thanks for contributing to FlowScope.

### Development quick start

```bash
docker compose up --build -d clickhouse seed
docker compose up --build -d
make test
```

### PR rules

1. Clear title and explain what changed and why.
2. Include test evidence.
3. Update docs when behavior/config changes.
4. Do not commit secrets or private data.
