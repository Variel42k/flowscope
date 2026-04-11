# Web Interface / Web интерфейс

FlowScope web UI is designed for operations and defensive monitoring workflows.  
Web UI FlowScope ориентирован на эксплуатацию и defensive-мониторинг.

## Main Pages / Основные страницы

1. **Overview**
- KPI cards: throughput, flows, active exporters
- top talkers, top protocols, top interfaces, top destination ports
- anomaly panel with heuristic detections

2. **Flows**
- active and historical flow tables
- filters by time, IP/subnet, exporter, protocol, port, interface, ASN, country
- row details drawer with full normalized record

3. **Sankey**
- source -> protocol/service -> destination traffic view
- click on a band to pivot into filtered flows

4. **Interaction Map**
- host/subnet/service interaction graph
- directional edges, zoom/pan, node search
- grouping by subnet, service, environment, ASN
- ego-network isolation and edge drill-down

5. **Alerts**
- rules and events view
- rule toggling and manual evaluation trigger

6. **Saved Views**
- create and reuse query presets for repeatable investigations

## Recommended Analyst Flow / Рекомендуемый сценарий аналитика

1. Start in **Overview** to detect suspicious or unusual traffic behavior.  
Начните с **Overview**, чтобы найти аномалии.
2. Pivot to **Interaction Map** for topology context.  
Перейдите в **Interaction Map** для топологического контекста.
3. Open edge details and inspect protocols/ports distribution.  
Откройте детали ребра и посмотрите распределение протоколов/портов.
4. Drill down into raw flows for confirmation.  
Провалитесь в сырые потоки для подтверждения.
