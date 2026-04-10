# Архитектура FlowScope

## Топология сервисов

- `collector` принимает UDP flow-телеметрию (NetFlow v5/v9, IPFIX, sFlow), декодирует, нормализует, обогащает и пишет в `raw_flow_events`.
- `worker` выполняет скользящие rollup-агрегации в `flow_1m_rollup`, `flow_1h_rollup`, `edges_1m`, `edges_1h` и обновляет `nodes_dimension`.
- `api` отдаёт REST endpoints для таблиц, графиков, Sankey и drill-down по interaction map.
- `web` (React UI) использует API и рендерит ECharts + Cytoscape визуализации.

## Пайплайн данных

1. UDP пакет принимается слушателем нужного протокола.
2. Декодер преобразует записи в канонический `FlowRecord`.
3. `normalize.Apply` вычисляет производные поля: subnet/private flags/time buckets/protocol labels/flow hash.
4. Enricher-цепочка добавляет метаданные (inventory, GeoIP/ASN, reverse DNS, interface aliases).
5. Raw flow пишется в ClickHouse.
6. Worker строит minute/hour rollups и edge/node измерения для быстрых запросов.
7. API использует rollups для дашбордов/карты и raw-таблицу для drill-down.

## Стратегия запросов

- Быстрые overview/topology запросы читают `edges_*` и `flow_*_rollup`.
- Drill-down (`/api/flows/*`, `/api/map/edge/:id`) читает `raw_flow_events`.
- Фильтры: время, IP/subnet, exporter, protocol, ports, interfaces, ASN, country, thresholds.

## Модель безопасности (MVP)

- Локальная auth-модель с одним admin пользователем из env vars.
- `POST /api/auth/login` возвращает JWT.
- Защищённые маршруты принимают Bearer JWT.
- Есть режим отключения auth для локальной отладки.

## Defensive-only scope

FlowScope реализует только observability и анализ сетевых взаимодействий.
Оффенсивные функции (сканирование, эксплуатация уязвимостей, persistence, credential theft, remote access tooling) не входят в проект.
