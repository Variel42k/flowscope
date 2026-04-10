# API Reference / Справочник API

## Русский

Базовый URL:
- `http://localhost:8088`
- через web proxy: `http://localhost:5173/api`

### Auth
- `POST /api/auth/login`

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

### Visuals
- `GET /api/sankey`
- `GET /api/map/graph`
- `GET /api/map/node/:id`
- `GET /api/map/edge/:id`

### Utility
- `GET /api/search`
- `POST /api/inventory/import`

Поддерживаемые query params: `from`, `to`, `page`, `page_size`, `sort_by`, `sort_dir`, сетевые и protocol filters.

## English

Base URLs:
- `http://localhost:8088`
- via web proxy: `http://localhost:5173/api`

### Auth
- `POST /api/auth/login`

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

### Visuals
- `GET /api/sankey`
- `GET /api/map/graph`
- `GET /api/map/node/:id`
- `GET /api/map/edge/:id`

### Utility
- `GET /api/search`
- `POST /api/inventory/import`

Supported query params include `from`, `to`, `page`, `page_size`, `sort_by`, `sort_dir`, and network/protocol filters.
