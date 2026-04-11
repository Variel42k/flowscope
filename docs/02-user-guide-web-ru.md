# Web UI User Guide / Руководство пользователя Web UI

## Русский

### Навигация

- `Overview` — сводка по трафику, топы, аномальные взаимодействия
- `Overview` — сводка по трафику, включая топ destination ports
- `Flows` — активные/исторические потоки + детали строки
- `Sankey` — агрегированная визуализация `source -> protocol/service -> destination`
- `Interaction Map` — интерактивный граф с drill-down
- `Alerts` — правила эвристик и список срабатываний
- `Saved Views` — реестр сохранённых фильтров

### Общие фильтры

`From`, `To`, `Search`, `Exporter`, `Protocol`, `Subnet`, `Source IP`, `Destination IP`, `Source Port`, `Destination Port`, `ASN`, `Country`.

На страницах `Overview/Flows/Sankey/Interaction Map` доступен блок `Saved Views` для сохранения и применения наборов фильтров (private/shared).

### Рекомендуемый сценарий расследования

1. Открыть `Overview` и выделить аномальное направление.
2. Перейти в `Interaction Map`, изолировать узел/ребро.
3. Проверить `Edge Details`.
4. Подтвердить события в `Matching Flows`.

## English

### Navigation

- `Overview` — traffic summary, top lists, suspicious interactions
- `Overview` — traffic summary including top destination ports
- `Flows` — active/historical flow tables with row details
- `Sankey` — aggregated `source -> protocol/service -> destination`
- `Interaction Map` — interactive graph with drill-down
- `Alerts` — heuristic rule management and firing history
- `Saved Views` — saved filter preset registry

### Global filters

`From`, `To`, `Search`, `Exporter`, `Protocol`, `Subnet`, `Source IP`, `Destination IP`, `Source Port`, `Destination Port`, `ASN`, `Country`.

`Saved Views` blocks on `Overview/Flows/Sankey/Interaction Map` let you save and reapply private/shared filter presets.

### Recommended investigation workflow

1. Start at `Overview` and identify unusual edges.
2. Open `Interaction Map` and isolate node/edge context.
3. Inspect `Edge Details`.
4. Confirm with `Matching Flows` raw entries.
