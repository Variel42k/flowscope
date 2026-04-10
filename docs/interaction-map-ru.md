# Interaction Map: поведение и возможности

## Узлы (Nodes)

Узлы могут представлять:

- host
- subnet
- service
- environment
- ASN группу

в зависимости от выбранного `mode` и `group_by`.

Каждый узел содержит:

- label и type
- bytes/packets in/out
- flow count и last seen
- теги (environment, service, country, ASN)
- private/internal флаги
- `collapsed` и `children` при группировке

## Рёбра (Edges)

Рёбра отражают наблюдаемое направленное взаимодействие:

- source/destination
- bytes, packets, flows
- список протоколов и top destination ports
- first seen / last seen
- exporter count

## Модель взаимодействия в UI

- Клик по узлу: подсветка соседей + загрузка node detail панели.
- Клик по ребру: edge detail панель + таблица matching flows.
- Group/collapse: серверная группировка через `group_by` (`subnet|service|environment|asn`).
- Ego isolate: `node_id` возвращает только подграф вокруг выбранного узла.
- Threshold filters: `min_bytes`, `min_flows` отсекают шумовые рёбра.
- Export: выделенный граф можно выгрузить в JSON.

## Рекомендации по анализу

1. Начать с `host_to_host` и небольшого временного окна.
2. Поднять пороги `min_bytes/min_flows`, чтобы убрать слабые связи.
3. Перейти к `internal_to_external` для поиска внешних коммуникаций.
4. Использовать edge drill-down для подтверждения на уровне raw flows.
