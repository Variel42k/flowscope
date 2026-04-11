# API Reference / Справочник API

Base URL / Базовый URL: `http://localhost:8088`

## Auth / Аутентификация

- `POST /api/auth/login`
- `GET /api/auth/me`
- `GET /api/auth/oidc/start` (optional)
- `GET /api/auth/oidc/callback` (optional)

## Core Data / Основные данные

- `GET /api/health`
- `GET /api/exporters`
- `GET /api/flows/active`
- `GET /api/flows/historical`
- `GET /api/flows/:id`

## Aggregations / Агрегации

- `GET /api/talkers/top`
- `GET /api/protocols/top`
- `GET /api/interfaces/top`
- `GET /api/ports/top`
- `GET /api/sankey`

## Map and Drill-down / Карта и детализация

- `GET /api/map/graph`
- `GET /api/map/node/:id`
- `GET /api/map/edge/:id`

## Alerts and Views / Алерты и представления

- `GET /api/alerts/rules`
- `POST /api/alerts/rules`
- `PUT /api/alerts/rules/:id`
- `DELETE /api/alerts/rules/:id`
- `GET /api/alerts/events`
- `POST /api/alerts/evaluate`
- `GET /api/views`
- `POST /api/views`
- `PUT /api/views/:id`
- `DELETE /api/views/:id`

## Utility / Служебные

- `GET /api/search`
- `POST /api/inventory/import`

## Common Query Parameters / Общие query-параметры

- `from`, `to`
- `page`, `page_size`
- `sort_by`, `sort_dir`
- filters: `ip`, `subnet`, `protocol`, `port`, `exporter`, `interface`, `asn`, `country`
