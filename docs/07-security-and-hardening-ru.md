# Безопасность и Hardening (RU)

## Минимум перед production

1. Сменить `FLOWSCOPE_ADMIN_USER`, `FLOWSCOPE_ADMIN_PASSWORD`, `FLOWSCOPE_JWT_SECRET`.
2. Убрать внешний доступ к ClickHouse (`8123`, `9000`).
3. Ограничить доступ к UDP listeners через firewall/ACL.
4. Запускать behind reverse proxy с TLS.
5. Не использовать `FLOWSCOPE_AUTH_DISABLED` вне dev.

## Сетевой hardening

- Разделить сети для:
  - ingest plane
  - API/Web plane
  - data plane (ClickHouse)
- Ограничить ingress только нужными источниками.

## Secrets management

- Хранить секреты в `.env` вне git.
- Предпочтительно использовать secret manager.
- Ротировать секреты на регулярной основе.

## Dependency hygiene

- Периодически обновлять Go/Node зависимости.
- Добавить автоматическую проверку уязвимостей в CI.
- Фиксировать критичные версии base image.

## Логирование и аудит

- Вести центральный сбор логов `api/collector/worker`.
- Контролировать аномальные пики login attempts.
- Мониторить рост объёма raw_flow_events и rollup latency.
