# API Reference (RU)

Базовый URL:

- локально: `http://localhost:8088`
- через web nginx proxy: `http://localhost:5173/api`

## Auth

### `POST /api/auth/login`

Request:

```json
{
  "username": "admin",
  "password": "admin123"
}
```

Response:

```json
{
  "token": "<jwt>",
  "user": "admin"
}
```

## Health

### `GET /api/health`

Response включает:

- `status`
- `service`
- `uptime_sec`
- `auth_enabled`

## Exporters

### `GET /api/exporters`

Возвращает paginated список экспортёров.

## Flows

### `GET /api/flows/active`
### `GET /api/flows/historical`
### `GET /api/flows/:id`

Используется для raw flow drill-down.

## Aggregates

### `GET /api/talkers/top`
### `GET /api/protocols/top`
### `GET /api/interfaces/top`

Параметр `limit` опционален.

## Visual APIs

### `GET /api/sankey`

Возвращает:

- `nodes[]`
- `links[]` (source/target/value/flows)

### `GET /api/map/graph`

Параметры:

- `mode`: `host_to_host|subnet_to_subnet|service_to_service|internal_to_external`
- `group_by`: `none|subnet|service|environment|asn`
- `min_bytes`
- `min_flows`
- `node_id` (ego isolation)
- плюс базовые фильтры по времени/поиску/экспортёру/протоколу

### `GET /api/map/node/:id`
### `GET /api/map/edge/:id`

Возвращает детали узла/ребра, series, top lists, matching flows.

## Search and Inventory

### `GET /api/search`

Query param:

- `search`

### `POST /api/inventory/import`

Request:

```json
{
  "items": [
    {
      "asset_id": "srv-1",
      "cidr": "10.10.1.10/32",
      "hostname": "app-01",
      "service": "frontend",
      "environment": "prod",
      "owner_team": "platform"
    }
  ]
}
```

## Common query params

Для list endpoints поддерживаются:

- `from`, `to` (RFC3339)
- `page`, `page_size`
- `sort_by`, `sort_dir`
- `src_ip`, `dst_ip`, `subnet`
- `exporter`
- `protocol`
- `src_port`, `dst_port`
- `input_interface`, `output_interface`
- `asn`, `country`
- `search`
- `min_bytes`, `min_flows`
