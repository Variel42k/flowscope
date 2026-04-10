package storage

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/flowscope/flowscope/internal/model"
)

func (r *Repository) Graph(ctx context.Context, f model.QueryFilter) (model.GraphResponse, error) {
	table := "edges_1m"
	if f.To.Sub(f.From) > 24*time.Hour {
		table = "edges_1h"
	}
	timeCol := "minute_bucket"
	if table == "edges_1h" {
		timeCol = "hour_bucket"
	}
	srcExpr, dstExpr := "src_node_id", "dst_node_id"
	switch f.Mode {
	case "subnet_to_subnet":
		srcExpr, dstExpr = "src_subnet", "dst_subnet"
	case "service_to_service":
		srcExpr, dstExpr = "if(src_service != '', src_service, src_subnet)", "if(dst_service != '', dst_service, dst_subnet)"
	case "internal_to_external":
		srcExpr, dstExpr = "src_node_id", "dst_node_id"
	}
	clauses := []string{timeCol + " >= ?", timeCol + " <= ?"}
	args := []any{f.From, f.To}
	if f.ExporterID != "" {
		clauses = append(clauses, "exporter_id = ?")
		args = append(args, f.ExporterID)
	}
	if f.Protocol != "" {
		clauses = append(clauses, "protocol = ?")
		args = append(args, strings.ToUpper(f.Protocol))
	}
	if f.MinBytes > 0 {
		clauses = append(clauses, "bytes >= ?")
		args = append(args, f.MinBytes)
	}
	if f.MinFlows > 0 {
		clauses = append(clauses, "flows >= ?")
		args = append(args, f.MinFlows)
	}
	if f.Mode == "internal_to_external" {
		clauses = append(clauses, "src_private != dst_private")
	}
	if f.Search != "" {
		needle := "%" + strings.ToLower(f.Search) + "%"
		clauses = append(clauses, "(lower(src_node_id) LIKE ? OR lower(dst_node_id) LIKE ? OR lower(src_service) LIKE ? OR lower(dst_service) LIKE ? OR lower(src_environment) LIKE ? OR lower(dst_environment) LIKE ?)")
		args = append(args, needle, needle, needle, needle, needle, needle)
	}
	where := strings.Join(clauses, " AND ")
	query := fmt.Sprintf(`
SELECT %s AS source,
       %s AS destination,
       sum(bytes) AS total_bytes,
       sum(packets) AS total_packets,
       sum(flows) AS total_flows,
       groupUniqArray(protocol) AS protocols,
       groupUniqArray(toUInt16(dst_port)) AS ports,
       min(first_seen) AS first_seen,
       max(last_seen) AS last_seen,
       uniqExact(exporter_id) AS exporters,
       any(src_private) AS src_private,
       any(dst_private) AS dst_private,
       any(src_environment) AS src_env,
       any(dst_environment) AS dst_env,
       any(src_country) AS src_country,
       any(dst_country) AS dst_country,
       any(src_asn) AS src_asn,
       any(dst_asn) AS dst_asn
FROM %s
WHERE %s
GROUP BY source, destination
ORDER BY total_bytes DESC
LIMIT 1500`, srcExpr, dstExpr, table, where)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return model.GraphResponse{}, err
	}
	defer rows.Close()
	nodes := map[string]*model.GraphNode{}
	edges := make([]model.GraphEdge, 0, 1500)
	for rows.Next() {
		var source, destination string
		var bytes, packets, flows, exporters uint64
		var protocols []string
		var ports []uint16
		var firstSeen, lastSeen time.Time
		var srcPrivate, dstPrivate bool
		var srcEnv, dstEnv, srcCountry, dstCountry string
		var srcASN, dstASN uint32
		if err := rows.Scan(&source, &destination, &bytes, &packets, &flows, &protocols, &ports, &firstSeen, &lastSeen, &exporters,
			&srcPrivate, &dstPrivate, &srcEnv, &dstEnv, &srcCountry, &dstCountry, &srcASN, &dstASN); err != nil {
			return model.GraphResponse{}, err
		}
		edges = append(edges, model.GraphEdge{
			ID: source + "->" + destination, Source: source, Destination: destination,
			Bytes: bytes, Packets: packets, Flows: flows,
			Protocols: protocols, TopDstPorts: ports,
			FirstSeen: firstSeen, LastSeen: lastSeen, Exporters: exporters, Direction: "directed",
		})
		if _, ok := nodes[source]; !ok {
			nodes[source] = &model.GraphNode{ID: source, Label: source, Type: inferNodeType(source), Tags: map[string]string{"environment": srcEnv, "country": srcCountry, "asn": strconv.FormatUint(uint64(srcASN), 10)}, Private: srcPrivate, Internal: srcPrivate}
		}
		if _, ok := nodes[destination]; !ok {
			nodes[destination] = &model.GraphNode{ID: destination, Label: destination, Type: inferNodeType(destination), Tags: map[string]string{"environment": dstEnv, "country": dstCountry, "asn": strconv.FormatUint(uint64(dstASN), 10)}, Private: dstPrivate, Internal: dstPrivate}
		}
		srcNode := nodes[source]
		dstNode := nodes[destination]
		srcNode.BytesOut += bytes
		srcNode.PacketsOut += packets
		srcNode.FlowCount += flows
		if lastSeen.After(srcNode.LastSeen) {
			srcNode.LastSeen = lastSeen
		}
		dstNode.BytesIn += bytes
		dstNode.PacketsIn += packets
		dstNode.FlowCount += flows
		if lastSeen.After(dstNode.LastSeen) {
			dstNode.LastSeen = lastSeen
		}
	}
	nodeList := make([]model.GraphNode, 0, len(nodes))
	for _, n := range nodes {
		nodeList = append(nodeList, *n)
	}
	sort.Slice(nodeList, func(i, j int) bool { return nodeList[i].FlowCount > nodeList[j].FlowCount })
	suspicious, _ := r.SuspiciousInteractions(ctx, f)
	return model.GraphResponse{Nodes: nodeList, Edges: edges, Mode: f.Mode, GroupBy: f.GroupBy, Suspicious: suspicious}, nil
}

func inferNodeType(s string) string {
	if strings.Contains(s, "/") {
		return "subnet"
	}
	if strings.Contains(s, ".") || strings.Contains(s, ":") {
		return "host"
	}
	return "service"
}

func (r *Repository) NodeDetails(ctx context.Context, nodeID string, f model.QueryFilter) (model.NodeDetails, error) {
	if f.From.IsZero() {
		f.From = time.Now().UTC().Add(-1 * time.Hour)
	}
	if f.To.IsZero() {
		f.To = time.Now().UTC()
	}
	baseWhere, baseArgs := buildFlowWhere(f, "")
	where := baseWhere + " AND (src_ip = ? OR dst_ip = ? OR src_hostname = ? OR dst_hostname = ? OR src_service = ? OR dst_service = ?)"
	args := append(baseArgs, nodeID, nodeID, nodeID, nodeID, nodeID, nodeID)
	out := model.NodeDetails{Node: model.GraphNode{ID: nodeID, Label: nodeID, Type: inferNodeType(nodeID)}}
	row := r.db.QueryRowContext(ctx, `
SELECT sumIf(bytes, dst_ip = ?), sumIf(packets, dst_ip = ?), countIf(dst_ip = ?),
       sumIf(bytes, src_ip = ?), sumIf(packets, src_ip = ?), countIf(src_ip = ?),
       max(timestamp_end)
FROM raw_flow_events
WHERE `+where, append([]any{nodeID, nodeID, nodeID, nodeID, nodeID, nodeID}, args...)...)
	if err := row.Scan(&out.Inbound.Bytes, &out.Inbound.Packets, &out.Inbound.Flows, &out.Outbound.Bytes, &out.Outbound.Packets, &out.Outbound.Flows, &out.Node.LastSeen); err != nil {
		return out, err
	}
	peerArgs := append([]any{nodeID}, baseArgs...)
	peerArgs = append(peerArgs, nodeID)
	peerArgs = append(peerArgs, baseArgs...)
	peerRows, err := r.db.QueryContext(ctx, `
SELECT peer, sum(bytes), sum(packets), count()
FROM (
	SELECT dst_ip AS peer, bytes, packets FROM raw_flow_events WHERE src_ip = ? AND `+baseWhere+`
	UNION ALL
	SELECT src_ip AS peer, bytes, packets FROM raw_flow_events WHERE dst_ip = ? AND `+baseWhere+`
)
GROUP BY peer
ORDER BY sum(bytes) DESC
LIMIT 10`, peerArgs...)
	if err == nil {
		defer peerRows.Close()
		out.TopPeers, _ = scanTopItems(peerRows)
	}
	commonArgs := append([]any{nodeID, nodeID}, baseArgs...)
	portsRows, err := r.db.QueryContext(ctx, `
SELECT toString(dst_port), sum(bytes), sum(packets), count()
FROM raw_flow_events
WHERE (src_ip = ? OR dst_ip = ?) AND `+baseWhere+`
GROUP BY dst_port
ORDER BY sum(bytes) DESC
LIMIT 10`, commonArgs...)
	if err == nil {
		defer portsRows.Close()
		out.TopPorts, _ = scanTopItems(portsRows)
	}
	protoRows, err := r.db.QueryContext(ctx, `
SELECT l4_protocol_name, sum(bytes), sum(packets), count()
FROM raw_flow_events
WHERE (src_ip = ? OR dst_ip = ?) AND `+baseWhere+`
GROUP BY l4_protocol_name
ORDER BY sum(bytes) DESC
LIMIT 10`, commonArgs...)
	if err == nil {
		defer protoRows.Close()
		out.TopProtocols, _ = scanTopItems(protoRows)
	}
	seriesRows, err := r.db.QueryContext(ctx, `
SELECT minute_bucket, sum(bytes), sum(packets), count()
FROM raw_flow_events
WHERE (src_ip = ? OR dst_ip = ?) AND `+baseWhere+`
GROUP BY minute_bucket
ORDER BY minute_bucket ASC`, commonArgs...)
	if err == nil {
		defer seriesRows.Close()
		for seriesRows.Next() {
			var p model.TimePoint
			if seriesRows.Scan(&p.Timestamp, &p.Bytes, &p.Packets, &p.Flows) == nil {
				out.Series = append(out.Series, p)
			}
		}
	}
	return out, nil
}

func (r *Repository) EdgeDetails(ctx context.Context, edgeID string, f model.QueryFilter) (model.EdgeDetails, error) {
	parts := strings.Split(edgeID, "->")
	if len(parts) != 2 {
		return model.EdgeDetails{}, fmt.Errorf("invalid edge id")
	}
	src, dst := parts[0], parts[1]
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 50
	}
	if f.From.IsZero() {
		f.From = time.Now().UTC().Add(-1 * time.Hour)
	}
	if f.To.IsZero() {
		f.To = time.Now().UTC()
	}
	where, args := buildFlowWhere(f, "")
	where = where + " AND src_ip = ? AND dst_ip = ?"
	args = append(args, src, dst)
	out := model.EdgeDetails{Edge: model.GraphEdge{ID: edgeID, Source: src, Destination: dst}}
	row := r.db.QueryRowContext(ctx, `
SELECT sum(bytes), sum(packets), count(), min(timestamp_start), max(timestamp_end), groupUniqArray(l4_protocol_name), groupUniqArray(dst_port), uniqExact(exporter_id)
FROM raw_flow_events
WHERE `+where, args...)
	if err := row.Scan(&out.Edge.Bytes, &out.Edge.Packets, &out.Edge.Flows, &out.Edge.FirstSeen, &out.Edge.LastSeen, &out.Edge.Protocols, &out.Edge.TopDstPorts, &out.Edge.Exporters); err != nil {
		return out, err
	}
	seriesRows, err := r.db.QueryContext(ctx, `
SELECT minute_bucket, sum(bytes), sum(packets), count()
FROM raw_flow_events
WHERE `+where+`
GROUP BY minute_bucket
ORDER BY minute_bucket`, args...)
	if err == nil {
		defer seriesRows.Close()
		for seriesRows.Next() {
			var p model.TimePoint
			if seriesRows.Scan(&p.Timestamp, &p.Bytes, &p.Packets, &p.Flows) == nil {
				out.Series = append(out.Series, p)
			}
		}
	}
	protoRows, err := r.db.QueryContext(ctx, `
SELECT l4_protocol_name, sum(bytes), sum(packets), count()
FROM raw_flow_events
WHERE `+where+`
GROUP BY l4_protocol_name
ORDER BY sum(bytes) DESC`, args...)
	if err == nil {
		defer protoRows.Close()
		out.Protocols, _ = scanTopItems(protoRows)
	}
	portRows, err := r.db.QueryContext(ctx, `
SELECT toString(dst_port), sum(bytes), sum(packets), count()
FROM raw_flow_events
WHERE `+where+`
GROUP BY dst_port
ORDER BY sum(bytes) DESC`, args...)
	if err == nil {
		defer portRows.Close()
		out.PortDist, _ = scanTopItems(portRows)
	}
	expRows, err := r.db.QueryContext(ctx, `
SELECT exporter_id, sum(bytes), sum(packets), count()
FROM raw_flow_events
WHERE `+where+`
GROUP BY exporter_id
ORDER BY sum(bytes) DESC`, args...)
	if err == nil {
		defer expRows.Close()
		out.Exporters, _ = scanTopItems(expRows)
	}
	offset := (f.Page - 1) * f.PageSize
	var flowTotal uint64
	if err := r.db.QueryRowContext(ctx, `SELECT count() FROM raw_flow_events WHERE `+where, args...).Scan(&flowTotal); err != nil {
		return out, err
	}
	flowQuery := fmt.Sprintf(`SELECT %s FROM raw_flow_events WHERE %s ORDER BY timestamp_end DESC LIMIT ? OFFSET ?`, flowSelectColumns(), where)
	flowRows, err := r.db.QueryContext(ctx, flowQuery, append(args, f.PageSize+1, offset)...)
	if err == nil {
		defer flowRows.Close()
		flows, scanErr := scanFlowRows(flowRows)
		if scanErr == nil {
			hasNext := len(flows) > f.PageSize
			if hasNext {
				flows = flows[:f.PageSize]
			}
			out.FlowPage = model.PageResult[model.FlowRecord]{
				Data:     flows,
				Page:     f.Page,
				PageSize: f.PageSize,
				Total:    flowTotal,
				HasNext:  hasNext,
				SortBy:   "timestamp_end",
				SortDir:  "DESC",
			}
		}
	}
	return out, nil
}

func (r *Repository) SuspiciousInteractions(ctx context.Context, f model.QueryFilter) ([]model.SuspiciousInteraction, error) {
	if f.From.IsZero() {
		f.From = time.Now().UTC().Add(-6 * time.Hour)
	}
	if f.To.IsZero() {
		f.To = time.Now().UTC()
	}
	out := make([]model.SuspiciousInteraction, 0, 20)
	rows, err := r.db.QueryContext(ctx, `
SELECT concat(src_node_id, '->', dst_node_id) AS edge_id, sum(bytes) AS total_bytes
FROM edges_1m
WHERE minute_bucket >= ? AND minute_bucket <= ?
GROUP BY edge_id
HAVING total_bytes > (
	SELECT quantileExact(0.99)(bytes) FROM edges_1m WHERE minute_bucket >= ? AND minute_bucket <= ?
)
ORDER BY total_bytes DESC
LIMIT 5`, f.From, f.To, f.From, f.To)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var edge string
			var b uint64
			if rows.Scan(&edge, &b) == nil {
				out = append(out, model.SuspiciousInteraction{Kind: "high_byte_edge", Description: fmt.Sprintf("Edge %s transferred %d bytes (outlier)", edge, b), EdgeID: edge, Severity: "medium"})
			}
		}
	}
	rows2, err := r.db.QueryContext(ctx, `
SELECT src_node_id, uniqExact(dst_node_id) AS fanout
FROM edges_1m
WHERE minute_bucket >= ? AND minute_bucket <= ? AND src_private = 1 AND dst_private = 0
GROUP BY src_node_id
HAVING fanout > 15
ORDER BY fanout DESC
LIMIT 5`, f.From, f.To)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var node string
			var fanout uint64
			if rows2.Scan(&node, &fanout) == nil {
				out = append(out, model.SuspiciousInteraction{Kind: "internal_to_many_external", Description: fmt.Sprintf("Internal node %s contacted %d external destinations", node, fanout), NodeID: node, Severity: "high"})
			}
		}
	}
	rows3, err := r.db.QueryContext(ctx, `
SELECT concat(src_node_id, '->', dst_node_id) AS edge_id, min(first_seen)
FROM edges_1m
WHERE minute_bucket >= ? AND minute_bucket <= ?
GROUP BY edge_id
HAVING min(first_seen) >= ?
ORDER BY min(first_seen) DESC
LIMIT 5`, f.From, f.To, f.From)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var edge string
			var first time.Time
			if rows3.Scan(&edge, &first) == nil {
				out = append(out, model.SuspiciousInteraction{Kind: "new_edge", Description: fmt.Sprintf("New edge observed: %s at %s", edge, first.Format(time.RFC3339)), EdgeID: edge, Severity: "low"})
			}
		}
	}
	rows4, err := r.db.QueryContext(ctx, `
SELECT src_node_id, dst_port, count() AS c
FROM edges_1m
WHERE minute_bucket >= ? AND minute_bucket <= ?
GROUP BY src_node_id, dst_port
HAVING c = 1
ORDER BY c ASC
LIMIT 5`, f.From, f.To)
	if err == nil {
		defer rows4.Close()
		for rows4.Next() {
			var node string
			var port uint16
			var c uint64
			if rows4.Scan(&node, &port, &c) == nil {
				out = append(out, model.SuspiciousInteraction{Kind: "port_outlier", Description: fmt.Sprintf("Node %s used unusual destination port %d", node, port), NodeID: node, Severity: "low"})
			}
		}
	}
	return out, nil
}
