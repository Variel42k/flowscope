# Web UI User Guide / Руководство пользователя Web UI

## Русский

### Навигация

- `Overview` — сводка по трафику, топы, аномальные взаимодействия
- `Flows` — активные/исторические потоки + детали строки
- `Sankey` — агрегированная визуализация `source -> protocol/service -> destination`
- `Interaction Map` — интерактивный граф с drill-down

### Общие фильтры

`From`, `To`, `Search`, `Exporter`, `Protocol`, `Subnet`, `Source IP`, `Destination IP`, `Source Port`, `Destination Port`, `ASN`, `Country`.

### Рекомендуемый сценарий расследования

1. Открыть `Overview` и выделить аномальное направление.
2. Перейти в `Interaction Map`, изолировать узел/ребро.
3. Проверить `Edge Details`.
4. Подтвердить события в `Matching Flows`.

## English

### Navigation

- `Overview` — traffic summary, top lists, suspicious interactions
- `Flows` — active/historical flow tables with row details
- `Sankey` — aggregated `source -> protocol/service -> destination`
- `Interaction Map` — interactive graph with drill-down

### Global filters

`From`, `To`, `Search`, `Exporter`, `Protocol`, `Subnet`, `Source IP`, `Destination IP`, `Source Port`, `Destination Port`, `ASN`, `Country`.

### Recommended investigation workflow

1. Start at `Overview` and identify unusual edges.
2. Open `Interaction Map` and isolate node/edge context.
3. Inspect `Edge Details`.
4. Confirm with `Matching Flows` raw entries.
