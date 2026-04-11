# Web Interface / Web интерфейс

## Русский

Web UI FlowScope ориентирован на эксплуатацию и defensive-мониторинг.

### Основные страницы

1. `Overview`
- KPI-карточки: throughput, flows, active exporters
- top talkers, protocols, interfaces, destination ports
- панель аномалий с эвристическими срабатываниями

2. `Flows`
- таблицы active/historical
- фильтры по времени, IP/subnet, exporter, protocol, port, interface, ASN, country
- drawer с деталями нормализованной записи

3. `Sankey`
- визуализация `source -> protocol/service -> destination`
- клик по ленте переводит в детализированную выборку flows

4. `Interaction Map`
- граф взаимодействий host/subnet/service
- directional edges, zoom/pan, поиск узлов
- группировка по subnet, service, environment, ASN
- ego-network isolation и edge drill-down

5. `Alerts`
- правила и события
- включение/выключение правил и ручной запуск evaluation

6. `Saved Views`
- сохранение и повторное применение фильтров

### Рекомендуемый сценарий аналитика

1. Начать с `Overview` и найти подозрительные направления.
2. Перейти в `Interaction Map` для контекста топологии.
3. Открыть edge details и проверить распределение protocols/ports.
4. Подтвердить гипотезу в raw flows.

## English

FlowScope web UI is designed for operations and defensive monitoring workflows.

### Main Pages

1. `Overview`
- KPI cards: throughput, flows, active exporters
- top talkers, protocols, interfaces, destination ports
- anomaly panel with heuristic detections

2. `Flows`
- active and historical tables
- filters by time, IP/subnet, exporter, protocol, port, interface, ASN, country
- details drawer with full normalized record

3. `Sankey`
- `source -> protocol/service -> destination` traffic view
- clicking a band pivots into filtered flows

4. `Interaction Map`
- host/subnet/service interaction graph
- directional edges, zoom/pan, node search
- grouping by subnet, service, environment, ASN
- ego-network isolation and edge drill-down

5. `Alerts`
- rules and events view
- rule toggling and manual evaluation trigger

6. `Saved Views`
- save and reuse filter presets

### Recommended Analyst Flow

1. Start in `Overview` to detect suspicious behavior.
2. Pivot to `Interaction Map` for topology context.
3. Open edge details and inspect protocols/ports distribution.
4. Drill down into raw flows for confirmation.
