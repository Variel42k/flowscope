# Capability Matrix / Карточка возможностей

## Русский

Эта страница показывает, что FlowScope уже умеет, что находится на уровне beta и что запланировано на следующие итерации.

### Статусы

- `Ready` — функциональность готова к практическому использованию в небольших и средних внедрениях.
- `Beta` — работает в MVP, но требует дополнительного hardening для широкого продакшн-использования.
- `Roadmap` — запланировано, но пока не реализовано.

### Матрица возможностей

| Область | Возможность | Статус | Примечания |
|---|---|---|---|
| Ingestion | NetFlow v5/v9, IPFIX, sFlow UDP listeners | Ready | Поддержка нескольких экспортеров, декодирование на базе parser-библиотек |
| Ingestion | Трекинг метаданных экспортера | Ready | First/last seen, source type, observation domain |
| Data Model | Каноническая схема нормализованного потока | Ready | Единая модель для всех входных протоколов |
| Enrichment | Inventory mapping (IP/subnet -> service/env/team) | Ready | Поддерживается импорт static YAML/JSON |
| Enrichment | Interface alias mapping | Ready | Наложение человекочитаемых имен интерфейсов из файла |
| Enrichment | GeoIP и rDNS | Beta | Опционально, зависит от локальных данных |
| Storage | ClickHouse raw flow events + rollups `1m/1h` | Ready | Time-partitioned модель с быстрыми агрегациями |
| Storage | Материализация node/edge для graph-запросов | Ready | Graph endpoints по умолчанию не сканируют raw-таблицу |
| API | Flows/exporters/search/talkers/protocols/interfaces/ports | Ready | Поддержка filter + pagination + sorting |
| API | Node/edge drill-down endpoints для interaction map | Ready | Протоколы, порты и matching flow page |
| API | Alert rules/events + manual evaluation endpoint | Ready | Эвристический pipeline интегрирован в worker |
| API | Saved views CRUD | Ready | Личные и shared presets с RBAC-проверками |
| Security | Local auth (JWT) + RBAC (`admin`/`viewer`) | Ready | Write-endpoints доступны только администратору |
| Security | OIDC SSO | Beta | Опционально, через переменные окружения |
| UI | Overview dashboard | Ready | KPI, топы, anomaly panel |
| UI | Flows explorer | Ready | Active/historical режимы и detail panel |
| UI | Sankey traffic visualization | Ready | Click-to-drill в таблицу потоков |
| UI | Interaction map (Cytoscape) | Ready | Grouping, ego isolate, threshold filters, export JSON |
| UI | Port-focused visualization | Ready | Top ports на Overview и port distribution в Edge Details |
| Operations | Docker Compose all-in-one stack | Ready | clickhouse + api + collector + worker + web + seed + flowgen |
| Operations | Seed dataset + synthetic traffic generation | Ready | Быстрый локальный демо-старт |
| Testing | Frontend lint/unit + backend package tests | Ready | Базовый CI-ready набор |
| Testing | End-to-end integration suite | Beta | Частично покрыто, можно расширять автоматизацию |
| Roadmap | Multi-tenant strict isolation | Roadmap | Пока не реализовано |
| Roadmap | Продвинутые anomaly-модели | Roadmap | Сейчас используются rule-based эвристики |
| Roadmap | HA-топология и autoscaling playbooks | Roadmap | Текущий таргет развертывания — compose/local-small prod |

### Практическое покрытие задач

#### SOC/NOC Visibility

- Кто с кем взаимодействует (host/subnet/service/internal-external)
- Объем трафика и количество потоков во времени
- Распределение портов и протоколов с edge-level drill-down

#### Engineering Operations

- Быстрый root-cause triage: от графа к raw flows
- Видимость здоровья экспортеров и continuity данных
- Повторяемые расследования через saved views

#### Security Monitoring (Defensive)

- Обнаружение новых communication edges
- Fan-out внутреннего хоста на большое число внешних адресов
- High-byte outlier edges
- Port outliers относительно базового профиля хоста

### Границы функционала

FlowScope — defensive observability-платформа. Проект намеренно **не** включает offensive-функции: сканирование, эксплуатацию уязвимостей, persistence-механики и инструменты удаленного доступа.

## English

This page summarizes what FlowScope can do now, what is in beta-level maturity, and what is planned for future iterations.

### Status Legend

- `Ready` — production-usable in small and medium deployments.
- `Beta` — works in MVP, but needs hardening for broader production use.
- `Roadmap` — planned, not yet implemented.

### Capability Matrix

| Domain | Capability | Status | Notes |
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
| Security | OIDC SSO | Beta | Optional integration via environment configuration |
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

### Practical Coverage

#### SOC/NOC Visibility

- Who talks to whom (host/subnet/service/internal-external)
- Traffic volume and flow count over time
- Port/protocol distribution with edge-level drill-down

#### Engineering Operations

- Fast root-cause triage from graph to raw flows
- Exporter health and data continuity visibility
- Repeatable investigations via saved filter views

#### Security Monitoring (Defensive)

- New communication edge detection
- Internal host fan-out to many external destinations
- High-byte edge outliers
- Port outliers against host baseline

### Scope Boundaries

FlowScope is a defensive observability platform. It intentionally does **not** include offensive capabilities such as scanning, exploitation, persistence tooling, or remote access frameworks.
