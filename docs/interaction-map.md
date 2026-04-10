# Interaction Map Behaviors

## Nodes

Nodes can represent hosts, subnets, services, environments, or ASN groups depending on selected graph mode/grouping.

Each node carries:

- label and type
- bytes/packets in/out
- flow count and last seen
- tags (environment, service, country, ASN)
- private/internal flags
- `collapsed` + `children` metadata when grouped

## Edges

Edges represent observed directional communication and include:

- source/destination
- bytes, packets, flow count
- protocol set and top destination ports
- first/last seen
- exporter count

## Interaction Model

- Click node: highlights ego neighborhood and requests node detail panel.
- Click edge: opens edge detail panel with flow drill-down table.
- Group/collapse: server-side grouping via `group_by` (subnet/service/environment/asn).
- Ego isolate: send `node_id` query param to return neighborhood-only subgraph.
- Threshold filtering: `min_bytes` and `min_flows` prune weak edges.
- Export: selected nodes/edges can be exported as JSON from UI.
