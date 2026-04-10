# FlowScope

[![License: AGPL v3+](https://img.shields.io/badge/license-AGPLv3%2B-blue.svg)](LICENSE)

FlowScope is an open-source, self-hosted **network flow observability platform** for defensive monitoring teams.

It ingests **NetFlow v5/v9, IPFIX, and sFlow**, stores normalized telemetry in ClickHouse, and provides a fast web UI for:

- top talkers, protocols, interfaces, and exporters
- active and historical flow exploration
- Sankey traffic analysis
- interactive node-edge interaction map with drill-down
- suspicious interaction heuristics (observability context)

License: **AGPL-3.0-or-later** (SPDX)

## Развёртывание (Docker Compose)

Ниже минимальная инструкция запуска, которая должна быть видна сразу на главной странице GitHub.

### 1. Требования

- Docker + Docker Compose v2
- свободные порты: `5173`, `8088`, `18123`, `19000`, `2055/udp`, `2056/udp`, `4739/udp`, `6343/udp`

### 2. Запуск

```bash
docker compose up --build -d
```

### 3. Проверка

```bash
docker compose ps
```

Ожидаемо:

- `clickhouse` в состоянии `healthy`
- `api`, `collector`, `worker`, `flowgen`, `web` в состоянии `Up`
- `seed` завершился с `Exited (0)`

### 4. Вход в Web UI

- Web: `http://localhost:5173`
- API: `http://localhost:8088`
- логин по умолчанию:
  - username: `admin`
  - password: `admin123`

Проверка API health:

```bash
curl http://localhost:8088/api/health
```

### 5. Остановка

```bash
docker compose down
```

С очисткой томов:

```bash
docker compose down -v
```

## Defensive-Only Scope

FlowScope is built for observability and blue-team operations only.
It does **not** implement offensive capabilities like scanning, packet injection, exploitation, persistence, credential theft, or remote access tooling.

## Why FlowScope

- Unified ingestion pipeline for multiple flow protocols
- Rollup-first querying strategy for responsive dashboards/graph views
- Deep drill-down path: Overview -> Node/Edge -> Flows
- Docker Compose friendly for local and small production deployments
- Clean Go + React monorepo, easy to extend

## Core Capabilities

1. UDP ingestion of NetFlow v5/v9, IPFIX, sFlow
2. Canonical flow normalization and optional enrichment
3. ClickHouse raw + rollup schema for recent and historical analysis
4. REST API for dashboards, map graph, search, and drill-down
5. Web UI with ECharts, TanStack Table, Cytoscape.js interaction map

## Stack

- Backend services: Go
- API: REST (`chi`)
- Storage: ClickHouse
- Frontend: React + TypeScript + Vite + Tailwind
- Charts: Apache ECharts
- Graph: Cytoscape.js
- Tables: TanStack Table
- Deployment: Docker Compose

## Repository Layout

```text
flowscope/
  cmd/
    api/
    collector/
    worker/
    seed/
    flowgen/
  internal/
    api/
    auth/
    config/
    decoder/
    enrich/
    graph/
    inventory/
    model/
    normalize/
    rollup/
    storage/
    util/
  deploy/
    clickhouse/init/
    docker/
    flowgen/
    seed/
  web/
  docs/
```

## Quick Start (EN)

```bash
docker compose up --build -d
```

Service endpoints:

- Web UI: `http://localhost:5173`
- API: `http://localhost:8088`
- ClickHouse HTTP: `http://localhost:18123`
- ClickHouse native: `localhost:19000`

Flow listeners:

- NetFlow v5: `localhost:2055/udp`
- NetFlow v9: `localhost:2056/udp`
- IPFIX: `localhost:4739/udp`
- sFlow: `localhost:6343/udp`

Default local demo credentials:

- username: `admin`
- password: `admin123`

Logs:

```bash
docker compose logs -f --tail=200
```

Stop and cleanup:

```bash
docker compose down -v
```

## API Surface (MVP)

- `POST /api/auth/login`
- `GET /api/health`
- `GET /api/exporters`
- `GET /api/flows/active`
- `GET /api/flows/historical`
- `GET /api/flows/:id`
- `GET /api/talkers/top`
- `GET /api/protocols/top`
- `GET /api/interfaces/top`
- `GET /api/sankey`
- `GET /api/map/graph`
- `GET /api/map/node/:id`
- `GET /api/map/edge/:id`
- `GET /api/search`
- `POST /api/inventory/import`

List/query endpoints support time range, filtering, pagination, and sorting.

## Documentation

- [Documentation index](docs/index.md)
- [Quick start (RU)](docs/01-quickstart-ru.md)
- [Architecture notes (EN)](docs/architecture.md)
- [Архитектура (RU)](docs/architecture-ru.md)
- [Interaction map details (EN)](docs/interaction-map.md)
- [Карта взаимодействий (RU)](docs/interaction-map-ru.md)
- [GitHub presentation kit (RU)](docs/github-presentation-ru.md)
- [Licensing guide (RU)](docs/10-licensing-ru.md)

## Security Notes

For local demos, defaults are intentionally simple.
Before production use, rotate credentials/secrets, restrict exposed ports, and apply hardening from [SECURITY.md](SECURITY.md).

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a PR.

## License

This project is licensed under **GNU AGPL v3.0 or later**.  
SPDX identifier: `AGPL-3.0-or-later`  
See [LICENSE](LICENSE) and [NOTICE](NOTICE).
