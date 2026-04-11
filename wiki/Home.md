# FlowScope Wiki

FlowScope is an open-source, self-hosted network flow observability platform for defensive teams.  
FlowScope — open-source self-hosted платформа наблюдаемости сетевых потоков для defensive-команд.

## Start in 5 Minutes / Запуск за 5 минут

1. Install Docker and Docker Compose v2.  
Установите Docker и Docker Compose v2.
2. Run:

```bash
docker compose up --build -d
```

3. Open:
- Web UI: `http://localhost:5173`
- API health: `http://localhost:8088/api/health`

Default credentials / Демо-учетные данные:
- login: `admin`
- password: `admin123`

## Core Capabilities / Основные возможности

- NetFlow v5/v9, IPFIX, sFlow ingestion over UDP
- ClickHouse raw storage + rollups (`1m`, `1h`)
- Traffic analytics: top talkers, protocols, interfaces, ports
- Sankey view and interaction map with drill-down to raw flows
- Alert rules/events, saved views, RBAC, optional OIDC

## Wiki Pages / Разделы Wiki

- [Quick Start / Быстрый старт](Quick-Start)
- [Web Interface / Web интерфейс](Web-Interface)
- [Architecture / Архитектура](Architecture)
- [API Reference / Справочник API](API-Reference)
- [Security and Hardening / Безопасность и hardening](Security-and-Hardening)
- [Operations / Эксплуатация](Operations)

## Source Docs / Документация в репозитории

- [Repository README](https://github.com/Variel42k/flowscope#readme)
- [Docs index](https://github.com/Variel42k/flowscope/blob/main/docs/index.md)
- [Capability matrix](https://github.com/Variel42k/flowscope/blob/main/docs/capability-matrix.md)
