# Capability Matrix / Карточка возможностей

This page summarizes what FlowScope can do right now, what is in beta-level maturity, and what is planned for next iterations.  
Эта страница показывает, что FlowScope уже умеет, что находится на уровне beta, и что запланировано в следующих итерациях.

## Status Legend / Обозначения статусов

- `Ready` — production-usable in small and medium deployments
- `Beta` — works in MVP, but needs hardening for broader production use
- `Roadmap` — planned, not yet implemented in current repository state

## Capability Matrix / Матрица возможностей

| Domain / Область | Capability / Возможность | Status | Notes / Примечания |
|---|---|---|---|
| Ingestion | NetFlow v5/v9, IPFIX, sFlow UDP listeners | Ready | Multi-exporter intake, parser-based decoding |
| Ingestion | Exporter metadata tracking | Ready | First/last seen, source type, observation domain |
| Data Model | Canonical normalized flow schema | Ready | Shared model across all input protocols |
| Enrichment | Inventory mapping (IP/subnet -> service/env/team) | Ready | Static YAML/JSON import supported |
| Enrichment | Interface alias mapping | Ready | File-based interface label overlay |
| Enrichment | GeoIP and rDNS | Beta | Optional and config-driven, depends on local data |
| Storage | ClickHouse raw flow events + 1m/1h rollups | Ready | Time-partitioned, aggregation-first query model |
| Storage | Node/edge materialization for graph queries | Ready | Graph endpoints avoid raw scans by default |
| API | Flows, exporters, search, talkers/protocols/interfaces/ports | Ready | Filter + pagination + sorting support |
| API | Interaction map node/edge drill-down endpoints | Ready | Includes protocols, ports, matching flow page |
| API | Alert rules/events + manual evaluation endpoint | Ready | Heuristic pipeline integrated with worker |
| API | Saved views CRUD | Ready | Personal and shared presets with RBAC checks |
| Security | Local auth (JWT) + RBAC (`admin`/`viewer`) | Ready | Admin-protected write endpoints |
| Security | OIDC SSO | Beta | Optional integration via env configuration |
| UI | Overview dashboard | Ready | KPIs, top entities, anomaly panel |
| UI | Flows explorer | Ready | Active/historical tabs + detail panel |
| UI | Sankey traffic visualization | Ready | Click-to-drill into flows |
| UI | Interaction map (Cytoscape) | Ready | Grouping, ego isolate, threshold filters, export JSON |
| UI | Port-focused visualization | Ready | Top ports on overview + edge port distribution |
| Operations | Docker Compose all-in-one stack | Ready | clickhouse + api + collector + worker + web + seed + flowgen |
| Operations | Seed dataset and synthetic traffic generation | Ready | Fast local demo bootstrap |
| Testing | Frontend lint/unit + backend package tests | Ready | CI-friendly, extendable baseline |
| Testing | End-to-end integration test suite | Beta | Partial checks, full scenario automation can be extended |
| Roadmap | Multi-tenant strict isolation | Roadmap | Not yet implemented |
| Roadmap | Advanced anomaly models (baseline learning) | Roadmap | Current heuristics are rule-based |
| Roadmap | HA topology and autoscaling playbooks | Roadmap | Current deployment target is compose/local-small prod |

## Practical Coverage / Практическое покрытие задач

### SOC/NOC Visibility

- Who talks to whom (host, subnet, service, internal/external)
- Traffic volume and flow count over time
- Port and protocol distribution with edge-level drill-down

### Engineering Operations

- Fast root-cause triage from graph to raw flows
- Exporter health and data continuity visibility
- Repeatable investigations via saved filter views

### Security Monitoring (Defensive)

- New communication edge detection
- Internal host fan-out to many external destinations
- High-byte edge outliers
- Port outliers compared to recent baseline

## Scope Boundaries / Границы функционала

FlowScope is a defensive observability platform. It intentionally does **not** include offensive capabilities such as scanning, exploitation, persistence tooling, or remote access frameworks.  
FlowScope — defensive observability-платформа. Она намеренно **не** включает offensive-возможности: сканирование, эксплуатацию уязвимостей, persistence-механики или инструменты удаленного доступа.
