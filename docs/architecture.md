# FlowScope Architecture Notes

## Service Topology

- `collector` receives UDP flow telemetry (NetFlow v5/v9, IPFIX, sFlow), decodes, normalizes, enriches, and inserts into `raw_flow_events`.
- `worker` performs sliding-window rollups into `flow_1m_rollup`, `flow_1h_rollup`, `edges_1m`, `edges_1h`, and refreshes `nodes_dimension`.
- `api` serves query/filter endpoints for tables, charts, Sankey, and interaction map drill-down.
- `web` is a React UI that consumes REST APIs and renders ECharts + Cytoscape visualizations.

## Data Pipeline

1. UDP packet ingested by protocol listener.
2. Decoder wrapper converts protocol records into canonical `FlowRecord`.
3. `normalize.Apply` derives key fields: subnets, privacy flags, time buckets, protocol labels, flow hash.
4. Enrichers attach optional metadata (inventory, GeoIP/ASN, reverse DNS, interface aliases).
5. Raw flow inserted into ClickHouse.
6. Worker aggregates minute/hour summaries and edge/node dimensions for efficient map/dashboard queries.
7. API defaults dashboards/map to rollups and uses raw table for drill-down details.

## Query Strategy

- Fast overview and topology queries rely on `edges_*` and `flow_*_rollup`.
- Drill-down (`/api/flows/*`, `/api/map/edge/:id`) reads from `raw_flow_events` with filter predicates.
- Filtering supports time range plus network-centric attributes: IP/subnet, exporter, protocol, ports, interfaces, ASN, country, thresholds.

## Security Model (MVP)

- Local auth with single admin user from env vars.
- `POST /api/auth/login` returns JWT.
- API accepts Bearer JWT for protected routes.
- Optional auth-disabled mode for local debugging.

## Defensive-only Scope

FlowScope includes only observability and traffic-analysis capabilities. It does not implement scanning, packet injection, exploitation, persistence, credential collection, or remote access tooling.
