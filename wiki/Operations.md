# Operations / Эксплуатация

## Service Lifecycle / Жизненный цикл сервисов

Start:

```bash
docker compose up --build -d
```

Status:

```bash
docker compose ps
docker compose logs -f api collector worker web
```

Stop:

```bash
docker compose down
```

## Backup / Резервное копирование

- Back up ClickHouse volumes regularly.
- Keep backup retention policy documented.
- Test restore procedure at least monthly.

## Restore / Восстановление

1. Stop stack.
2. Restore ClickHouse data volumes from backup.
3. Start stack and verify health endpoint.

## Upgrade / Обновление

1. Pull latest repository changes.
2. Review `.env.example` for new settings.
3. Rebuild and restart:

```bash
docker compose up --build -d
```

4. Validate key flows in UI: Overview -> Map -> edge drill-down.
