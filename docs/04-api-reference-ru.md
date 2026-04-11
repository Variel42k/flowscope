# API Reference / Справочник API

## Русский

Базовый URL:
- `http://localhost:8088`
- через web proxy: `http://localhost:5173/api`

### Auth
- `POST /api/auth/login`
- `GET /api/auth/oidc/start`
- `GET /api/auth/oidc/callback`
- `GET /api/auth/me`

### Health
- `GET /api/health`

### Data
- `GET /api/exporters`
- `GET /api/flows/active`
- `GET /api/flows/historical`
- `GET /api/flows/:id`

### Aggregates
- `GET /api/talkers/top`
- `GET /api/protocols/top`
- `GET /api/interfaces/top`
- `GET /api/ports/top`

### Visuals
- `GET /api/sankey`
- `GET /api/map/graph`
- `GET /api/map/node/:id`
- `GET /api/map/edge/:id`

### Utility
- `GET /api/search`
- `POST /api/inventory/import`
- `GET /api/views`
- `POST /api/views`
- `PUT /api/views/:id`
- `DELETE /api/views/:id`
- `GET /api/alerts/rules`
- `POST /api/alerts/rules`
- `PUT /api/alerts/rules/:id`
- `DELETE /api/alerts/rules/:id`
- `GET /api/alerts/events`
- `POST /api/alerts/evaluate`

Поддерживаемые query params: `from`, `to`, `page`, `page_size`, `sort_by`, `sort_dir`, сетевые и protocol filters.

## English

Base URLs:
- `http://localhost:8088`
- via web proxy: `http://localhost:5173/api`

### Auth
- `POST /api/auth/login`
- `GET /api/auth/oidc/start`
- `GET /api/auth/oidc/callback`
- `GET /api/auth/me`

### Health
- `GET /api/health`

### Data
- `GET /api/exporters`
- `GET /api/flows/active`
- `GET /api/flows/historical`
- `GET /api/flows/:id`

### Aggregates
- `GET /api/talkers/top`
- `GET /api/protocols/top`
- `GET /api/interfaces/top`
- `GET /api/ports/top`

### Visuals
- `GET /api/sankey`
- `GET /api/map/graph`
- `GET /api/map/node/:id`
- `GET /api/map/edge/:id`

### Utility
- `GET /api/search`
- `POST /api/inventory/import`
- `GET /api/views`
- `POST /api/views`
- `PUT /api/views/:id`
- `DELETE /api/views/:id`
- `GET /api/alerts/rules`
- `POST /api/alerts/rules`
- `PUT /api/alerts/rules/:id`
- `DELETE /api/alerts/rules/:id`
- `GET /api/alerts/events`
- `POST /api/alerts/evaluate`

Supported query params include `from`, `to`, `page`, `page_size`, `sort_by`, `sort_dir`, and network/protocol filters.
