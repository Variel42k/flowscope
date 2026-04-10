# Interaction Map Behavior / Поведение interaction map

## English

Nodes represent hosts/subnets/services/environments/ASN groups depending on mode and grouping.
Edges represent directional communication with bytes, packets, flows, protocols, ports, first/last seen.

Main interactions:
- click node -> highlight neighbors + node details
- click edge -> edge details + matching flows
- group/collapse by subnet/service/environment/asn
- ego isolation via `node_id`
- threshold filters via `min_bytes` and `min_flows`
- export selected graph as JSON

## Русский

Узлы представляют хосты/подсети/сервисы/окружения/ASN группы в зависимости от режима.
Рёбра отражают направленные взаимодействия и содержат bytes, packets, flows, protocols, ports, first/last seen.

Основные действия:
- клик по узлу -> подсветка соседей + детали узла
- клик по ребру -> детали ребра + matching flows
- группировка/схлопывание по subnet/service/environment/asn
- изоляция ego-сети через `node_id`
- пороги через `min_bytes` и `min_flows`
- экспорт выбранного подграфа в JSON
