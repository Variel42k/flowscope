# Резервное копирование и восстановление (RU)

## Что бэкапить

Для Compose-поставки критичен volume `clickhouse_data`.

## Бэкап (базовый подход)

1. Остановить запись (опционально: остановить `collector/flowgen`).
2. Выполнить backup volume средствами вашей платформы.

Пример с остановкой сервисов:

```powershell
docker compose stop collector flowgen worker api
```

После backup:

```powershell
docker compose start api worker collector flowgen
```

## Логический экспорт (дополнительно)

Можно экспортировать таблицы через ClickHouse HTTP/native клиент в формате SQL/CSV/Parquet.

## Восстановление

1. Поднять чистый стенд.
2. Восстановить volume `clickhouse_data`.
3. Запустить сервисы:

```powershell
docker compose up -d
```

4. Проверить:
   - `GET /api/health`
   - наличие данных в `Overview` / `Flows`

## Рекомендуемый RPO/RTO для MVP

- RPO: 15-60 минут (зависит от частоты backup)
- RTO: 15-30 минут для small deployment
