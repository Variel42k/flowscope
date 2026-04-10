# Backup and Restore / Резервное копирование и восстановление

## Русский

Критичный объект backup: volume `clickhouse_data`.

Базовый сценарий:
1. Остановить запись (`collector`, `flowgen`, `worker`, при необходимости `api`).
2. Сделать backup volume.
3. Восстановить volume и поднять сервисы.

## English

Critical backup target: `clickhouse_data` volume.

Basic flow:
1. Stop write-heavy services (`collector`, `flowgen`, `worker`, optionally `api`).
2. Backup the volume.
3. Restore the volume and start services.
