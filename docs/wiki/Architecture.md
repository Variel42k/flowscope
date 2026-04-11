# Architecture / Архитектура

## Русский

### Компоненты верхнего уровня

- `collector`: принимает NetFlow/IPFIX/sFlow по UDP, нормализует и обогащает потоки
- `clickhouse`: хранение raw и rollup-данных
- `api`: REST backend, auth, search, drill-down queries
- `worker`: rollups, graph aggregates, alert evaluation
- `web`: React frontend для dashboard, Sankey, map и расследований
- `flowgen` + `seed`: генерация демо-трафика и bootstrap-данных

### Путь данных

1. Exporters отправляют телеметрию на UDP listeners.
2. Collector декодирует и нормализует записи в каноническую схему.
3. Данные пишутся в `raw_flow_events`.
4. Worker строит rollups (`flow_1m_rollup`, `flow_1h_rollup`, `edges_1m`, `edges_1h`).
5. API по умолчанию обслуживает dashboards и graph-query из rollups.
6. Raw-таблица используется для точечного drill-down.

### Границы доверия

- Недоверенная зона: UDP exporters и внешние API-клиенты
- Доверенная зона: внутренняя сеть сервисов и ClickHouse
- Auth boundary: login/JWT и проверка ролей
- Admin boundary: write-endpoints для inventory, alerts и конфигурационных операций

### Режим развертывания

- Основной режим: Docker Compose (local/dev/small production)
- Продакшн-рекомендация: reverse proxy + TLS + ограничение сетевой экспозиции

## English

### High-Level Components

- `collector`: receives NetFlow/IPFIX/sFlow over UDP, normalizes and enriches flows
- `clickhouse`: raw and rollup storage
- `api`: REST backend, auth, search, drill-down queries
- `worker`: rollups, graph aggregates, alert evaluation
- `web`: React frontend for dashboards, Sankey, map, and investigations
- `flowgen` + `seed`: demo traffic generation and deterministic bootstrap

### Data Path

1. Exporters send telemetry to UDP listeners.
2. Collector decodes and normalizes records to canonical schema.
3. Data is inserted into `raw_flow_events`.
4. Worker builds rollups (`flow_1m_rollup`, `flow_1h_rollup`, `edges_1m`, `edges_1h`).
5. API serves dashboards and graph queries from rollups by default.
6. Raw table is used for fine-grained drill-down.

### Trust Boundaries

- Untrusted input: UDP exporters and external API clients
- Trusted zone: internal service network and ClickHouse
- Auth boundary: login/JWT and role checks
- Admin boundary: write endpoints for inventory, alerts, and configuration operations

### Deployment Mode

- Primary mode: Docker Compose (local/dev/small production)
- Production recommendation: reverse proxy + TLS + restricted network exposure
