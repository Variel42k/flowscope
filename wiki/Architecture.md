# Architecture / Архитектура

## High-Level Components / Компоненты верхнего уровня

- `collector`: receives NetFlow/IPFIX/sFlow over UDP, normalizes and enriches flows
- `clickhouse`: raw and rollup storage for analytics
- `api`: REST backend, auth, search, drill-down queries
- `worker`: rollups, graph aggregates, alert evaluation
- `web`: React frontend for dashboards, Sankey, map, and investigations
- `flowgen` + `seed`: demo data generation and deterministic bootstrap

## Data Path / Путь данных

1. Exporters send flow telemetry to UDP listeners.
2. Collector decodes and normalizes to canonical schema.
3. Data is inserted into `raw_flow_events`.
4. Worker builds rollups (`flow_1m_rollup`, `flow_1h_rollup`, `edges_1m`, `edges_1h`).
5. API serves dashboards and graph queries from rollups by default.
6. Raw table is used for fine-grained drill-down.

## Security Boundaries / Границы доверия

- Untrusted input: UDP exporters and external API clients
- Trusted zone: internal service network and ClickHouse
- Auth boundary: API login/JWT and role checks
- Admin boundary: write endpoints for inventory, alerts, and configuration operations

## Deployment Mode / Режим развертывания

- Primary mode: Docker Compose (local/dev/small production)
- Production recommendation: reverse proxy + TLS + restricted network exposure
