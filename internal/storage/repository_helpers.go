package storage

import (
	"database/sql"

	"github.com/flowscope/flowscope/internal/model"
)

func flowSelectColumns() string {
	return `flow_id, timestamp_start, timestamp_end, exporter_id, exporter_ip, exporter_name, observation_id,
	src_ip, dst_ip, src_port, dst_port, ip_protocol, l4_protocol_name,
	bytes, packets, input_interface, output_interface, input_interface_alias, output_interface_alias,
	src_asn, dst_asn, src_country, dst_country,
	src_hostname, dst_hostname, src_service, dst_service,
	src_environment, dst_environment, src_owner_team, dst_owner_team,
	src_mac, dst_mac, vlan_id, tcp_flags, flow_direction, sampler_rate,
	src_subnet, dst_subnet, src_is_private, dst_is_private, flow_key_hash,
	minute_bucket, hour_bucket, source_type`
}

func scanFlowRows(rows *sql.Rows) ([]model.FlowRecord, error) {
	out := make([]model.FlowRecord, 0, 64)
	for rows.Next() {
		var rec model.FlowRecord
		if err := scanFlow(rows, &rec); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanFlow(s scanner, rec *model.FlowRecord) error {
	return s.Scan(
		&rec.FlowID, &rec.TimestampStart, &rec.TimestampEnd, &rec.ExporterID, &rec.ExporterIP, &rec.ExporterName, &rec.ObservationID,
		&rec.SrcIP, &rec.DstIP, &rec.SrcPort, &rec.DstPort, &rec.IPProtocol, &rec.L4ProtocolName,
		&rec.Bytes, &rec.Packets, &rec.InputInterface, &rec.OutputInterface, &rec.InputIfAlias, &rec.OutputIfAlias,
		&rec.SrcASN, &rec.DstASN, &rec.SrcCountry, &rec.DstCountry,
		&rec.SrcHostname, &rec.DstHostname, &rec.SrcService, &rec.DstService,
		&rec.SrcEnvironment, &rec.DstEnvironment, &rec.SrcOwnerTeam, &rec.DstOwnerTeam,
		&rec.SrcMAC, &rec.DstMAC, &rec.VLANID, &rec.TCPFlags, &rec.FlowDirection, &rec.SamplerRate,
		&rec.SrcSubnet, &rec.DstSubnet, &rec.SrcIsPrivate, &rec.DstIsPrivate, &rec.FlowKeyHash,
		&rec.MinuteBucket, &rec.HourBucket, &rec.SourceType,
	)
}

func scanTopItems(rows *sql.Rows) ([]model.TopItem, error) {
	items := make([]model.TopItem, 0, 32)
	for rows.Next() {
		var item model.TopItem
		if err := rows.Scan(&item.Key, &item.Bytes, &item.Packets, &item.Flows); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
