# Operations / Эксплуатация

## Русский

### Жизненный цикл сервисов

Запуск:

```bash
docker compose up --build -d
```

Статус:

```bash
docker compose ps
docker compose logs -f api collector worker web
```

Остановка:

```bash
docker compose down
```

### Резервное копирование

- Регулярно делайте backup ClickHouse volumes.
- Зафиксируйте политику хранения backup-копий.
- Проверяйте restore-процедуру минимум раз в месяц.

### Восстановление

1. Остановите stack.
2. Восстановите ClickHouse volumes из backup.
3. Запустите stack и проверьте health endpoint.

### Обновление

1. Получите последние изменения репозитория.
2. Проверьте `.env.example` на новые переменные.
3. Пересоберите и перезапустите:

```bash
docker compose up --build -d
```

4. Проверьте ключевой путь в UI: Overview -> Map -> edge drill-down.

## English

### Service Lifecycle

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

### Backup

- Back up ClickHouse volumes regularly.
- Document backup retention policy.
- Test restore procedure at least monthly.

### Restore

1. Stop the stack.
2. Restore ClickHouse volumes from backup.
3. Start the stack and verify the health endpoint.

### Upgrade

1. Pull latest repository changes.
2. Review `.env.example` for new settings.
3. Rebuild and restart:

```bash
docker compose up --build -d
```

4. Validate key UI flow: Overview -> Map -> edge drill-down.
