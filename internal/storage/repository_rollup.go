package storage

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) ApplyRetention(ctx context.Context, days int) error {
	if days <= 0 {
		days = 30
	}
	hourDays := days * 6
	queries := []string{
		fmt.Sprintf("ALTER TABLE raw_flow_events MODIFY TTL minute_bucket + INTERVAL %d DAY", days),
		fmt.Sprintf("ALTER TABLE flow_1m_rollup MODIFY TTL minute_bucket + INTERVAL %d DAY", days),
		fmt.Sprintf("ALTER TABLE edges_1m MODIFY TTL minute_bucket + INTERVAL %d DAY", days),
		fmt.Sprintf("ALTER TABLE flow_1h_rollup MODIFY TTL hour_bucket + INTERVAL %d DAY", hourDays),
		fmt.Sprintf("ALTER TABLE edges_1h MODIFY TTL hour_bucket + INTERVAL %d DAY", hourDays),
	}
	for _, q := range queries {
		if _, err := r.db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) RunRollups(ctx context.Context, from, to time.Time) error {
	queries := []struct {
		sql  string
		args []any
	}{
		{`ALTER TABLE flow_1m_rollup DELETE WHERE minute_bucket >= ? AND minute_bucket < ?`, []any{from, to}},
		{`INSERT INTO flow_1m_rollup
		 SELECT minute_bucket, exporter_id, src_subnet, dst_subnet, l4_protocol_name, dst_port,
		        input_interface, output_interface, sum(bytes), sum(packets), count()
		 FROM raw_flow_events
		 WHERE minute_bucket >= ? AND minute_bucket < ?
		 GROUP BY minute_bucket, exporter_id, src_subnet, dst_subnet, l4_protocol_name, dst_port, input_interface, output_interface`, []any{from, to}},
		{`ALTER TABLE edges_1m DELETE WHERE minute_bucket >= ? AND minute_bucket < ?`, []any{from, to}},
		{`INSERT INTO edges_1m
		 SELECT minute_bucket, exporter_id,
		        src_ip, dst_ip, src_subnet, dst_subnet,
		        src_service, dst_service, src_environment, dst_environment,
		        src_asn, dst_asn, src_country, dst_country,
		        src_is_private, dst_is_private, l4_protocol_name, dst_port,
		        sum(bytes), sum(packets), count(), min(timestamp_start), max(timestamp_end)
		 FROM raw_flow_events
		 WHERE minute_bucket >= ? AND minute_bucket < ?
		 GROUP BY minute_bucket, exporter_id, src_ip, dst_ip, src_subnet, dst_subnet,
		          src_service, dst_service, src_environment, dst_environment,
		          src_asn, dst_asn, src_country, dst_country,
		          src_is_private, dst_is_private, l4_protocol_name, dst_port`, []any{from, to}},
		{`ALTER TABLE flow_1h_rollup DELETE WHERE hour_bucket >= toStartOfHour(?) AND hour_bucket < toStartOfHour(?)`, []any{from, to}},
		{`INSERT INTO flow_1h_rollup
		 SELECT toStartOfHour(minute_bucket), exporter_id, src_subnet, dst_subnet, protocol, dst_port,
		        sum(bytes), sum(packets), sum(flows)
		 FROM edges_1m
		 WHERE minute_bucket >= ? AND minute_bucket < ?
		 GROUP BY toStartOfHour(minute_bucket), exporter_id, src_subnet, dst_subnet, protocol, dst_port`, []any{from, to}},
		{`ALTER TABLE edges_1h DELETE WHERE hour_bucket >= toStartOfHour(?) AND hour_bucket < toStartOfHour(?)`, []any{from, to}},
		{`INSERT INTO edges_1h
		 SELECT toStartOfHour(minute_bucket), exporter_id, src_node_id, dst_node_id,
		        src_subnet, dst_subnet, src_service, dst_service, src_environment, dst_environment,
		        src_asn, dst_asn, src_country, dst_country, src_private, dst_private, protocol, dst_port,
		        sum(bytes), sum(packets), sum(flows), min(first_seen), max(last_seen)
		 FROM edges_1m
		 WHERE minute_bucket >= ? AND minute_bucket < ?
		 GROUP BY toStartOfHour(minute_bucket), exporter_id, src_node_id, dst_node_id,
		          src_subnet, dst_subnet, src_service, dst_service, src_environment, dst_environment,
		          src_asn, dst_asn, src_country, dst_country, src_private, dst_private, protocol, dst_port`, []any{from, to}},
		{`ALTER TABLE nodes_dimension DELETE WHERE updated_at >= ? AND updated_at < ?`, []any{from, to}},
		{`INSERT INTO nodes_dimension
		 SELECT node_id, any(node_type), sum(bytes_in), sum(bytes_out), sum(packets_in), sum(packets_out), sum(flows), max(last_seen), any(tags_json), max(updated_at)
		 FROM (
		 	SELECT src_ip AS node_id, 'host' AS node_type, 0 AS bytes_in, sum(bytes) AS bytes_out, 0 AS packets_in, sum(packets) AS packets_out,
		 	       count() AS flows, max(timestamp_end) AS last_seen,
		 	       toJSONString(map('environment', any(src_environment), 'service', any(src_service), 'country', any(src_country), 'asn', toString(any(src_asn)))) AS tags_json,
		 	       max(timestamp_end) AS updated_at
		 	FROM raw_flow_events
		 	WHERE minute_bucket >= ? AND minute_bucket < ?
		 	GROUP BY src_ip
		 	UNION ALL
		 	SELECT dst_ip AS node_id, 'host' AS node_type, sum(bytes) AS bytes_in, 0 AS bytes_out, sum(packets) AS packets_in, 0 AS packets_out,
		 	       count() AS flows, max(timestamp_end) AS last_seen,
		 	       toJSONString(map('environment', any(dst_environment), 'service', any(dst_service), 'country', any(dst_country), 'asn', toString(any(dst_asn)))) AS tags_json,
		 	       max(timestamp_end) AS updated_at
		 	FROM raw_flow_events
		 	WHERE minute_bucket >= ? AND minute_bucket < ?
		 	GROUP BY dst_ip
		 )
		 GROUP BY node_id`, []any{from, to, from, to}},
	}
	for i, q := range queries {
		if _, err := r.db.ExecContext(ctx, q.sql, q.args...); err != nil {
			return fmt.Errorf("rollup query %d failed: %w", i, err)
		}
	}
	return nil
}
