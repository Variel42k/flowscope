# Руководство администратора и Runbook (RU)

## Базовые операции

Запуск:

```powershell
docker compose up --build -d
```

Статус:

```powershell
docker compose ps
```

Логи:

```powershell
docker compose logs -f --tail=200
```

Остановка:

```powershell
docker compose down
```

Полная очистка:

```powershell
docker compose down -v
```

## Компоненты

- `collector`: UDP ingestion и нормализация
- `api`: REST endpoints и auth
- `worker`: rollups/materialization
- `flowgen`: фоновый синтетический трафик для демо
- `seed`: импорт стартовых данных
- `clickhouse`: хранение raw и rollup
- `web`: UI

## Health checks

API:

```powershell
Invoke-RestMethod -Uri http://localhost:8088/api/health
```

Ожидаемо: `"status": "ok"`.

## Повторный seed

```powershell
docker compose run --rm seed
```

## Обновление конфигурации

1. Измените env values в `docker-compose.yml` или `.env`.
2. Пересоздайте затронутые сервисы:

```powershell
docker compose up -d --build --force-recreate api collector worker web
```

## Полезные операционные проверки

1. Убедиться, что worker регулярно обновляет rollups.
2. Убедиться, что collector принимает UDP пакеты и вставляет raw flows.
3. Проверить, что Web UI отображает новые данные.

## Минимальный production checklist

1. Заменить дефолтные креды и JWT secret.
2. Ограничить exposed ports и доступ к ClickHouse.
3. Включить TLS на ingress/reverse proxy.
4. Выключить `FLOWSCOPE_AUTH_DISABLED`.
5. Включить резервное копирование данных ClickHouse.
