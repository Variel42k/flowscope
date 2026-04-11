# FlowScope Wiki

## Русский

FlowScope — open-source self-hosted платформа наблюдаемости сетевых потоков для defensive-команд.

### Запуск за 5 минут

1. Установите Docker и Docker Compose v2.
2. Выполните:

```bash
docker compose up --build -d
```

3. Откройте:
- Web UI: `http://localhost:5173`
- API health: `http://localhost:8088/api/health`

Демо-учетные данные:
- login: `admin`
- password: `admin123`

### Основные возможности

- Прием NetFlow v5/v9, IPFIX, sFlow по UDP
- Хранение сырых потоков и rollup-агрегатов (`1m`, `1h`) в ClickHouse
- Аналитика трафика: top talkers, protocols, interfaces, ports
- Sankey и Interaction Map с переходом к raw flows
- Alert rules/events, saved views, RBAC, optional OIDC

### Разделы Wiki

- [Быстрый старт / Quick Start](Quick-Start)
- [Web интерфейс / Web Interface](Web-Interface)
- [Архитектура / Architecture](Architecture)
- [Справочник API / API Reference](API-Reference)
- [Безопасность / Security and Hardening](Security-and-Hardening)
- [Эксплуатация / Operations](Operations)

### Документация в репозитории

- [Repository README](https://github.com/Variel42k/flowscope#readme)
- [Docs index](https://github.com/Variel42k/flowscope/blob/main/docs/index.md)
- [Capability matrix](https://github.com/Variel42k/flowscope/blob/main/docs/capability-matrix.md)

## English

FlowScope is an open-source, self-hosted network flow observability platform for defensive teams.

### Start in 5 Minutes

1. Install Docker and Docker Compose v2.
2. Run:

```bash
docker compose up --build -d
```

3. Open:
- Web UI: `http://localhost:5173`
- API health: `http://localhost:8088/api/health`

Demo credentials:
- login: `admin`
- password: `admin123`

### Core Capabilities

- NetFlow v5/v9, IPFIX, sFlow UDP ingestion
- ClickHouse raw storage plus rollups (`1m`, `1h`)
- Traffic analytics: top talkers, protocols, interfaces, ports
- Sankey and Interaction Map with drill-down to raw flows
- Alert rules/events, saved views, RBAC, optional OIDC

### Wiki Pages

- [Quick Start](Quick-Start)
- [Web Interface](Web-Interface)
- [Architecture](Architecture)
- [API Reference](API-Reference)
- [Security and Hardening](Security-and-Hardening)
- [Operations](Operations)

### Source Docs

- [Repository README](https://github.com/Variel42k/flowscope#readme)
- [Docs index](https://github.com/Variel42k/flowscope/blob/main/docs/index.md)
- [Capability matrix](https://github.com/Variel42k/flowscope/blob/main/docs/capability-matrix.md)
