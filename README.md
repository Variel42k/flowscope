# FlowScope

[![License: AGPL v3+](https://img.shields.io/badge/license-AGPLv3%2B-blue.svg)](LICENSE)

FlowScope is an open-source, self-hosted **network flow observability platform** for defensive monitoring teams.  
FlowScope — open-source self-hosted платформа наблюдаемости сетевых потоков для defensive monitoring.

License: **AGPL-3.0-or-later** (SPDX)

## Deployment / Развёртывание

### Русский

1. Требования:
- Docker + Docker Compose v2
- свободные порты: `5173`, `8088`, `18123`, `19000`, `2055/udp`, `2056/udp`, `4739/udp`, `6343/udp`

2. Запуск:

```bash
docker compose up --build -d
```

3. Проверка:

```bash
docker compose ps
curl http://localhost:8088/api/health
```

4. Вход:
- Web: `http://localhost:5173`
- API: `http://localhost:8088`
- demo login: `admin` / `admin123`

5. Остановка:

```bash
docker compose down
# с очисткой данных:
docker compose down -v
```

### English

1. Requirements:
- Docker + Docker Compose v2
- open ports: `5173`, `8088`, `18123`, `19000`, `2055/udp`, `2056/udp`, `4739/udp`, `6343/udp`

2. Start:

```bash
docker compose up --build -d
```

3. Verify:

```bash
docker compose ps
curl http://localhost:8088/api/health
```

4. Access:
- Web: `http://localhost:5173`
- API: `http://localhost:8088`
- demo login: `admin` / `admin123`

5. Stop:

```bash
docker compose down
# wipe volumes:
docker compose down -v
```

## Core Features / Основные возможности

- NetFlow v5/v9, IPFIX, sFlow ingestion
- canonical flow normalization + optional enrichment
- ClickHouse raw + rollup analytics
- dashboards (talkers/protocols/interfaces/exporters)
- Sankey traffic view
- interaction map with node/edge drill-down

## API (MVP)

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

## Documentation / Документация

- [Docs index](docs/index.md)
- [Quick start](docs/01-quickstart-ru.md)
- [Web UI user guide](docs/02-user-guide-web-ru.md)
- [Admin operations](docs/03-admin-operations-ru.md)
- [API reference](docs/04-api-reference-ru.md)
- [Configuration reference](docs/05-configuration-reference-ru.md)
- [Troubleshooting](docs/06-troubleshooting-ru.md)
- [Security and hardening](docs/07-security-and-hardening-ru.md)
- [Backup and restore](docs/08-backup-restore-ru.md)
- [Development guide](docs/09-development-ru.md)
- [Licensing](docs/10-licensing-ru.md)
- [Architecture](docs/architecture.md)
- [Interaction map behavior](docs/interaction-map.md)
- [GitHub presentation kit](docs/github-presentation-ru.md)

## Security / Безопасность

Before production usage: rotate credentials/secrets, restrict exposed ports, and apply [SECURITY.md](SECURITY.md).  
Перед production: смените креды/секреты, ограничьте внешние порты и следуйте [SECURITY.md](SECURITY.md).

## Contributing / Вклад

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License / Лицензия

GNU AGPL v3.0 or later (`AGPL-3.0-or-later`). See [LICENSE](LICENSE) and [NOTICE](NOTICE).
