package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/flowscope/flowscope/internal/model"
)

type Repository struct {
	db            *sql.DB
	retentionDays int
}

func NewRepository(db *sql.DB, retentionDays int) *Repository {
	return &Repository{db: db, retentionDays: retentionDays}
}

func (r *Repository) Health(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *Repository) InsertFlows(ctx context.Context, records []model.FlowRecord) error {
	if len(records) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, `
INSERT INTO raw_flow_events (
	flow_id, timestamp_start, timestamp_end, exporter_id, exporter_ip, exporter_name, observation_id,
	src_ip, dst_ip, src_port, dst_port, ip_protocol, l4_protocol_name,
	bytes, packets, input_interface, output_interface, input_interface_alias, output_interface_alias,
	src_asn, dst_asn, src_country, dst_country,
	src_hostname, dst_hostname, src_service, dst_service,
	src_environment, dst_environment, src_owner_team, dst_owner_team,
	src_mac, dst_mac, vlan_id, tcp_flags, flow_direction, sampler_rate,
	src_subnet, dst_subnet, src_is_private, dst_is_private, flow_key_hash,
	minute_bucket, hour_bucket, source_type
) VALUES (
	?,?,?,?,?,?,?,
	?,?,?,?,?, ?,
	?,?,?,?, ?,?,?,?,
	?,?,?, ?,?,
	?,?,?,?, ?,?,?,?,
	?,?,?, ?,?,
	?,?,?, ?,?,
	?,?,?
)
`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, rec := range records {
		_, err := stmt.ExecContext(ctx,
			rec.FlowID, rec.TimestampStart, rec.TimestampEnd, rec.ExporterID, rec.ExporterIP, rec.ExporterName, rec.ObservationID,
			rec.SrcIP, rec.DstIP, rec.SrcPort, rec.DstPort, rec.IPProtocol, rec.L4ProtocolName,
			rec.Bytes, rec.Packets, rec.InputInterface, rec.OutputInterface, rec.InputIfAlias, rec.OutputIfAlias,
			rec.SrcASN, rec.DstASN, rec.SrcCountry, rec.DstCountry,
			rec.SrcHostname, rec.DstHostname, rec.SrcService, rec.DstService,
			rec.SrcEnvironment, rec.DstEnvironment, rec.SrcOwnerTeam, rec.DstOwnerTeam,
			rec.SrcMAC, rec.DstMAC, rec.VLANID, rec.TCPFlags, rec.FlowDirection, rec.SamplerRate,
			rec.SrcSubnet, rec.DstSubnet, rec.SrcIsPrivate, rec.DstIsPrivate, rec.FlowKeyHash,
			rec.MinuteBucket, rec.HourBucket, rec.SourceType,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) UpsertExporter(ctx context.Context, exp model.Exporter) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO exporters (exporter_id, exporter_ip, exporter_name, observation_id, first_seen, last_seen, flows_observed, last_source_type)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`, exp.ExporterID, exp.ExporterIP, exp.ExporterName, exp.ObservationID, exp.FirstSeen, exp.LastSeen, exp.FlowsObserved, exp.LastSourceType)
	return err
}

func (r *Repository) ListExporters(ctx context.Context, f model.QueryFilter) (model.PageResult[model.Exporter], error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 25
	}
	if f.PageSize > 200 {
		f.PageSize = 200
	}
	offset := (f.Page - 1) * f.PageSize
	var total uint64
	countRow := r.db.QueryRowContext(ctx, `SELECT uniqExact(exporter_id) FROM exporters`)
	if err := countRow.Scan(&total); err != nil {
		return model.PageResult[model.Exporter]{}, err
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT exporter_id, anyLast(exporter_ip), anyLast(exporter_name), anyLast(observation_id), min(first_seen), max(last_seen), sum(flows_observed), anyLast(last_source_type)
FROM exporters
GROUP BY exporter_id
ORDER BY max(last_seen) DESC
LIMIT ? OFFSET ?`, f.PageSize+1, offset)
	if err != nil {
		return model.PageResult[model.Exporter]{}, err
	}
	defer rows.Close()
	items := make([]model.Exporter, 0, f.PageSize)
	for rows.Next() {
		var it model.Exporter
		if err := rows.Scan(&it.ExporterID, &it.ExporterIP, &it.ExporterName, &it.ObservationID, &it.FirstSeen, &it.LastSeen, &it.FlowsObserved, &it.LastSourceType); err != nil {
			return model.PageResult[model.Exporter]{}, err
		}
		items = append(items, it)
	}
	hasNext := len(items) > f.PageSize
	if hasNext {
		items = items[:f.PageSize]
	}
	return model.PageResult[model.Exporter]{Data: items, Page: f.Page, PageSize: f.PageSize, Total: total, HasNext: hasNext, SortBy: "last_seen", SortDir: "DESC"}, nil
}

func (r *Repository) QueryActiveFlows(ctx context.Context, f model.QueryFilter) (model.PageResult[model.FlowRecord], error) {
	if f.To.IsZero() {
		f.To = time.Now().UTC()
	}
	if f.From.IsZero() {
		f.From = f.To.Add(-15 * time.Minute)
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 50
	}
	if f.PageSize > 500 {
		f.PageSize = 500
	}
	f.SortBy = safeSort(f.SortBy, []string{"timestamp_start", "timestamp_end", "bytes", "packets", "src_ip", "dst_ip"}, "timestamp_end")
	f.SortDir = safeSortDir(f.SortDir)
	where, args := buildFlowWhere(f, "")
	countQuery := `SELECT count() FROM raw_flow_events WHERE ` + where
	var total uint64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return model.PageResult[model.FlowRecord]{}, err
	}
	args = append(args, f.PageSize+1, (f.Page-1)*f.PageSize)
	query := fmt.Sprintf(`
SELECT %s
FROM raw_flow_events
WHERE %s
ORDER BY %s %s
LIMIT ? OFFSET ?`, flowSelectColumns(), where, f.SortBy, f.SortDir)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return model.PageResult[model.FlowRecord]{}, err
	}
	defer rows.Close()
	data, err := scanFlowRows(rows)
	if err != nil {
		return model.PageResult[model.FlowRecord]{}, err
	}
	hasNext := len(data) > f.PageSize
	if hasNext {
		data = data[:f.PageSize]
	}
	return model.PageResult[model.FlowRecord]{Data: data, Page: f.Page, PageSize: f.PageSize, Total: total, HasNext: hasNext, SortBy: f.SortBy, SortDir: f.SortDir}, nil
}

func (r *Repository) QueryHistoricalFlows(ctx context.Context, f model.QueryFilter) (model.PageResult[model.FlowRecord], error) {
	return r.QueryActiveFlows(ctx, f)
}

func (r *Repository) GetFlowByID(ctx context.Context, id string) (*model.FlowRecord, error) {
	query := fmt.Sprintf(`SELECT %s FROM raw_flow_events WHERE flow_id = ? LIMIT 1`, flowSelectColumns())
	row := r.db.QueryRowContext(ctx, query, id)
	var rec model.FlowRecord
	if err := scanFlow(row, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *Repository) TopTalkers(ctx context.Context, f model.QueryFilter, limit int) ([]model.TopItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	where, args := buildFlowWhere(f, "")
	query := `
SELECT ip, sum(bytes) AS bytes, sum(packets) AS packets, sum(flows) AS flows
FROM (
	SELECT src_ip AS ip, bytes, packets, 1 AS flows FROM raw_flow_events WHERE ` + where + `
	UNION ALL
	SELECT dst_ip AS ip, bytes, packets, 1 AS flows FROM raw_flow_events WHERE ` + where + `
)
GROUP BY ip
ORDER BY bytes DESC
LIMIT ?`
	args = append(args, args...)
	args = append(args, limit)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTopItems(rows)
}

func (r *Repository) TopProtocols(ctx context.Context, f model.QueryFilter, limit int) ([]model.TopItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	where, args := buildFlowWhere(f, "")
	args = append(args, limit)
	rows, err := r.db.QueryContext(ctx, `
SELECT l4_protocol_name AS key, sum(bytes) AS bytes, sum(packets) AS packets, count() AS flows
FROM raw_flow_events
WHERE `+where+`
GROUP BY l4_protocol_name
ORDER BY bytes DESC
LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTopItems(rows)
}

func (r *Repository) TopInterfaces(ctx context.Context, f model.QueryFilter, limit int) ([]model.TopItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	where, args := buildFlowWhere(f, "")
	query := `
SELECT if(alias != '', concat(toString(interface_id), ' (', alias, ')'), toString(interface_id)) AS key, sum(bytes) AS bytes, sum(packets) AS packets, count() AS flows
FROM (
	SELECT input_interface AS interface_id, input_interface_alias AS alias, bytes, packets FROM raw_flow_events WHERE ` + where + `
	UNION ALL
	SELECT output_interface AS interface_id, output_interface_alias AS alias, bytes, packets FROM raw_flow_events WHERE ` + where + `
)
GROUP BY interface_id, alias
ORDER BY bytes DESC
LIMIT ?`
	args = append(args, args...)
	args = append(args, limit)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTopItems(rows)
}

func (r *Repository) Sankey(ctx context.Context, f model.QueryFilter) (model.SankeyResponse, error) {
	where, args := buildFlowWhere(f, "")
	rows, err := r.db.QueryContext(ctx, `
SELECT src, mid, dst, sum(bytes) AS total_bytes, count() AS flows
FROM (
	SELECT if(src_hostname != '', src_hostname, src_ip) AS src,
	       if(dst_service != '', dst_service, l4_protocol_name) AS mid,
	       if(dst_hostname != '', dst_hostname, dst_ip) AS dst,
	       bytes
	FROM raw_flow_events
	WHERE `+where+`
)
GROUP BY src, mid, dst
ORDER BY total_bytes DESC
LIMIT 400`, args...)
	if err != nil {
		return model.SankeyResponse{}, err
	}
	defer rows.Close()
	nodes := map[string]struct{}{}
	links := make([]model.SankeyLink, 0, 800)
	for rows.Next() {
		var src, mid, dst string
		var bytes, flows uint64
		if err := rows.Scan(&src, &mid, &dst, &bytes, &flows); err != nil {
			return model.SankeyResponse{}, err
		}
		nodes[src], nodes[mid], nodes[dst] = struct{}{}, struct{}{}, struct{}{}
		links = append(links,
			model.SankeyLink{Source: src, Target: mid, Value: bytes, Flows: flows},
			model.SankeyLink{Source: mid, Target: dst, Value: bytes, Flows: flows},
		)
	}
	nodeList := make([]model.SankeyNode, 0, len(nodes))
	for name := range nodes {
		nodeList = append(nodeList, model.SankeyNode{Name: name})
	}
	return model.SankeyResponse{Nodes: nodeList, Links: links}, nil
}

func (r *Repository) Search(ctx context.Context, f model.QueryFilter) ([]model.TopItem, error) {
	needle := "%" + strings.ToLower(strings.TrimSpace(f.Search)) + "%"
	rows, err := r.db.QueryContext(ctx, `
SELECT key, bytes, packets, flows
FROM (
	SELECT src_ip AS key, sum(bytes) AS bytes, sum(packets) AS packets, count() AS flows FROM raw_flow_events WHERE lower(src_ip) LIKE ? GROUP BY src_ip
	UNION ALL
	SELECT dst_ip AS key, sum(bytes) AS bytes, sum(packets) AS packets, count() AS flows FROM raw_flow_events WHERE lower(dst_ip) LIKE ? GROUP BY dst_ip
	UNION ALL
	SELECT src_hostname AS key, sum(bytes) AS bytes, sum(packets) AS packets, count() AS flows FROM raw_flow_events WHERE lower(src_hostname) LIKE ? AND src_hostname != '' GROUP BY src_hostname
	UNION ALL
	SELECT dst_hostname AS key, sum(bytes) AS bytes, sum(packets) AS packets, count() AS flows FROM raw_flow_events WHERE lower(dst_hostname) LIKE ? AND dst_hostname != '' GROUP BY dst_hostname
	UNION ALL
	SELECT src_service AS key, sum(bytes) AS bytes, sum(packets) AS packets, count() AS flows FROM raw_flow_events WHERE lower(src_service) LIKE ? AND src_service != '' GROUP BY src_service
	UNION ALL
	SELECT dst_service AS key, sum(bytes) AS bytes, sum(packets) AS packets, count() AS flows FROM raw_flow_events WHERE lower(dst_service) LIKE ? AND dst_service != '' GROUP BY dst_service
)
ORDER BY bytes DESC
LIMIT 30`, needle, needle, needle, needle, needle, needle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTopItems(rows)
}

func (r *Repository) ImportInventory(ctx context.Context, items []model.InventoryAsset) error {
	if len(items) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, `
INSERT INTO inventory_assets (asset_id, cidr, hostname, service, environment, owner_team, interface_aliases, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	now := time.Now().UTC()
	for _, it := range items {
		aliases, _ := json.Marshal(it.InterfaceMap)
		if _, err := stmt.ExecContext(ctx, it.AssetID, it.CIDR, it.Hostname, it.Service, it.Environment, it.OwnerTeam, string(aliases), now); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
