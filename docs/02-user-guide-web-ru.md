# Руководство пользователя Web UI (RU)

## Навигация

- `Overview` — сводка по трафику, топы, аномальные взаимодействия.
- `Flows` — таблицы активных/исторических потоков с фильтрами и деталями.
- `Sankey` — потоковая визуализация `source -> protocol/service -> destination`.
- `Interaction Map` — интерактивный граф узлов и рёбер с drill-down.

## Общие фильтры

Доступны на всех страницах:

- `From` / `To` (временное окно)
- `Search` (IP/hostname/service)
- `Exporter`
- `Protocol`
- `Subnet`
- `Source IP` / `Destination IP`
- `Source Port` / `Destination Port`
- `ASN`
- `Country`

## Страница Overview

Показывает:

- Total Throughput
- Total Flows
- Active Exporters
- Top Edges
- Top Protocols/Top Talkers charts
- Recent Suspicious Interactions

Используйте её как отправную точку для расследования.

## Страница Flows

Режимы:

- `Active Flows`
- `Historical Flows`

Что можно делать:

1. Применить фильтры.
2. Выбрать строку в таблице.
3. Смотреть детали потока в правой панели:
   - IP/port pair
   - protocol
   - bytes/packets
   - exporter
   - interfaces
   - ASN/country

## Страница Sankey

Позволяет:

1. Выбрать временное окно и фильтры.
2. Кликнуть на “ленточную связь”.
3. Получить таблицу matching flows снизу.

Используйте для быстрого понимания “кто с кем и через что” в агрегированном виде.

## Страница Interaction Map

Основные функции:

- zoom/pan
- выбор `Graph Mode`:
  - `host_to_host`
  - `subnet_to_subnet`
  - `service_to_service`
  - `internal_to_external`
- `Group / Collapse`:
  - `none`
  - `subnet`
  - `service`
  - `environment`
  - `asn`
- пороги `Min Bytes`, `Min Flows`
- `Isolate Ego Node`
- `Export JSON`

Поведение:

1. Клик по node — подсветка соседей + панель node details.
2. Клик по edge — панель edge details + matching flows table.
3. Кнопка `Expand Selected Group` — быстрый выход из группировки в конкретный узел.

## Рекомендуемый workflow расследования

1. `Overview` -> заметить аномалию/edge.
2. `Interaction Map` -> изолировать узел или ребро.
3. `Edge Details` -> проверить порты/протоколы/экспортёры.
4. `Matching Flows` -> подтверждение на уровне raw flow events.
