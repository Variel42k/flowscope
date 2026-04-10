# GitHub Presentation Kit (RU)

Этот файл содержит готовые тексты для оформления репозитория FlowScope в GitHub.

## 1) About -> Description (до 160 символов)

`FlowScope is a self-hosted network flow observability platform with NetFlow/IPFIX/sFlow ingestion, ClickHouse analytics, Sankey and interaction map drill-down.`

## 2) About -> Website

Если пока нет публичного сайта, оставьте пустым.
Для демо можно указывать документацию репозитория:

`https://github.com/<owner>/<repo>/blob/main/docs/index.md`

## 3) About -> Topics

Рекомендуемые topics:

- `network-observability`
- `netflow`
- `ipfix`
- `sflow`
- `clickhouse`
- `golang`
- `react`
- `cytoscape`
- `echarts`
- `self-hosted`
- `soc`
- `noc`
- `defensive-security`

## 4) Короткий pitch (для README intro / соцсетей)

`FlowScope helps infrastructure and security teams understand network interactions in real time and historically, from exporter to edge to underlying flows.`

## 5) Release notes template (v0.1.0 MVP)

```markdown
## FlowScope v0.1.0 (MVP)

### Highlights
- UDP ingestion for NetFlow v5/v9, IPFIX, sFlow
- Canonical flow schema with optional enrichment
- ClickHouse raw + rollup storage model
- Web dashboards: top talkers/protocols/interfaces/exporters
- Sankey traffic exploration
- Interactive node-edge interaction map with drill-down

### API
- `/api/flows/*`, `/api/talkers/top`, `/api/protocols/top`, `/api/interfaces/top`
- `/api/sankey`, `/api/map/graph`, `/api/map/node/:id`, `/api/map/edge/:id`
- `/api/search`, `/api/inventory/import`

### Operations
- Docker Compose local deployment
- Seed loader + synthetic flow generator for demos

### Important
- Defensive monitoring use only
- Rotate default credentials and secrets before production usage
```

## 6) Рекомендуемый pinned text (для профиля/README)

`Open-source defensive network flow observability: ingest, store, query, visualize, and investigate traffic interactions from graph to raw flow records.`

## 7) Что добавить визуально в GitHub

1. Social preview image (`Settings -> Social preview`) с заголовком `FlowScope`.
2. Скриншоты страниц `Overview`, `Sankey`, `Interaction Map` в `docs/images`.
3. Значки статуса CI после настройки workflow.
