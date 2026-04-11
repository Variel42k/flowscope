# API Reference / Справочник API

Base URL / Базовый URL: `http://localhost:8088`

## Русский

### Auth

- `POST /api/auth/login`
- `GET /api/auth/me`
- `GET /api/auth/oidc/start` (optional)
- `GET /api/auth/oidc/callback` (optional)

### Основные данные

- `GET /api/health`
- `GET /api/exporters`
- `GET /api/flows/active`
- `GET /api/flows/historical`
- `GET /api/flows/:id`

### Агрегации

- `GET /api/talkers/top`
- `GET /api/protocols/top`
- `GET /api/interfaces/top`
- `GET /api/ports/top`
- `GET /api/sankey`

### Карта и drill-down

- `GET /api/map/graph`
- `GET /api/map/node/:id`
- `GET /api/map/edge/:id`

### Alerts и Views

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

### Служебные

- `GET /api/search`
- `POST /api/inventory/import`

### Общие query-параметры

- `from`, `to`
- `page`, `page_size`
- `sort_by`, `sort_dir`
- filters: `ip`, `subnet`, `protocol`, `port`, `exporter`, `interface`, `asn`, `country`

## English

### Auth

- `POST /api/auth/login`
- `GET /api/auth/me`
- `GET /api/auth/oidc/start` (optional)
- `GET /api/auth/oidc/callback` (optional)

### Core Data

- `GET /api/health`
- `GET /api/exporters`
- `GET /api/flows/active`
- `GET /api/flows/historical`
- `GET /api/flows/:id`

### Aggregations

- `GET /api/talkers/top`
- `GET /api/protocols/top`
- `GET /api/interfaces/top`
- `GET /api/ports/top`
- `GET /api/sankey`

### Map and Drill-down

- `GET /api/map/graph`
- `GET /api/map/node/:id`
- `GET /api/map/edge/:id`

### Alerts and Views

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

### Utility

- `GET /api/search`
- `POST /api/inventory/import`

### Common Query Parameters

- `from`, `to`
- `page`, `page_size`
- `sort_by`, `sort_dir`
- filters: `ip`, `subnet`, `protocol`, `port`, `exporter`, `interface`, `asn`, `country`
