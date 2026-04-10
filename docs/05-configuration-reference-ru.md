# Справочник конфигурации (RU)

## Основные переменные

- `FLOWSCOPE_HTTP_ADDR` — адрес API, по умолчанию `:8088`
- `FLOWSCOPE_CLICKHOUSE_DSN` — DSN ClickHouse
- `FLOWSCOPE_JWT_SECRET` — секрет подписи JWT
- `FLOWSCOPE_ADMIN_USER` — логин администратора
- `FLOWSCOPE_ADMIN_PASSWORD` — пароль администратора
- `FLOWSCOPE_AUTH_DISABLED` — отключение auth (только для dev)

## Ingestion

- `FLOWSCOPE_LISTEN_NETFLOW_V5` — UDP listener NetFlow v5
- `FLOWSCOPE_LISTEN_NETFLOW_V9` — UDP listener NetFlow v9
- `FLOWSCOPE_LISTEN_IPFIX` — UDP listener IPFIX
- `FLOWSCOPE_LISTEN_SFLOW` — UDP listener sFlow

## Rollups and retention

- `FLOWSCOPE_RETENTION_DAYS` — retention raw/1m
- `FLOWSCOPE_ROLLUP_INTERVAL` — период worker rollup
- `FLOWSCOPE_DEFAULT_QUERY_WINDOW` — дефолтное окно запросов

## Enrichment

- `FLOWSCOPE_ENABLE_INVENTORY`
- `FLOWSCOPE_INVENTORY_PATH`
- `FLOWSCOPE_INVENTORY_FORMAT` (`yaml|json`)
- `FLOWSCOPE_ENABLE_GEOIP`
- `FLOWSCOPE_GEOIP_CITY_PATH`
- `FLOWSCOPE_GEOIP_ASN_PATH`
- `FLOWSCOPE_ENABLE_RDNS`
- `FLOWSCOPE_RDNS_CACHE_TTL`
- `FLOWSCOPE_ENABLE_INTERFACE_ALIAS`
- `FLOWSCOPE_INTERFACE_ALIAS_PATH`

## Seed/generator

- `FLOWSCOPE_SEED_FIXTURE_PATH`
- `FLOWSCOPE_SEED_EXIT_SUCCESS`

## Рекомендации для production

1. Не использовать demo credentials.
2. Хранить секреты вне compose-файлов.
3. Ограничить доступ к портам ClickHouse и UDP listeners.
4. Включить TLS и reverse proxy security headers.
