package storage

import (
	"fmt"
	"strings"

	"github.com/flowscope/flowscope/internal/model"
)

func buildFlowWhere(f model.QueryFilter, tableAlias string) (string, []any) {
	alias := ""
	if tableAlias != "" {
		alias = tableAlias + "."
	}
	clauses := []string{fmt.Sprintf("%stimestamp_start >= ?", alias), fmt.Sprintf("%stimestamp_start <= ?", alias)}
	args := []any{f.From, f.To}

	if f.SrcIP != "" {
		clauses = append(clauses, fmt.Sprintf("%ssrc_ip = ?", alias))
		args = append(args, f.SrcIP)
	}
	if f.DstIP != "" {
		clauses = append(clauses, fmt.Sprintf("%sdst_ip = ?", alias))
		args = append(args, f.DstIP)
	}
	if f.Subnet != "" {
		clauses = append(clauses, fmt.Sprintf("(%ssrc_subnet = ? OR %sdst_subnet = ?)", alias, alias))
		args = append(args, f.Subnet, f.Subnet)
	}
	if f.ExporterID != "" {
		clauses = append(clauses, fmt.Sprintf("%sexporter_id = ?", alias))
		args = append(args, f.ExporterID)
	}
	if f.Protocol != "" {
		clauses = append(clauses, fmt.Sprintf("%sl4_protocol_name = ?", alias))
		args = append(args, strings.ToUpper(f.Protocol))
	}
	if f.SrcPort > 0 {
		clauses = append(clauses, fmt.Sprintf("%ssrc_port = ?", alias))
		args = append(args, f.SrcPort)
	}
	if f.DstPort > 0 {
		clauses = append(clauses, fmt.Sprintf("%sdst_port = ?", alias))
		args = append(args, f.DstPort)
	}
	if f.InputInterface > 0 {
		clauses = append(clauses, fmt.Sprintf("%sinput_interface = ?", alias))
		args = append(args, f.InputInterface)
	}
	if f.OutputInterface > 0 {
		clauses = append(clauses, fmt.Sprintf("%soutput_interface = ?", alias))
		args = append(args, f.OutputInterface)
	}
	if f.ASN > 0 {
		clauses = append(clauses, fmt.Sprintf("(%ssrc_asn = ? OR %sdst_asn = ?)", alias, alias))
		args = append(args, f.ASN, f.ASN)
	}
	if f.Country != "" {
		clauses = append(clauses, fmt.Sprintf("(%ssrc_country = ? OR %sdst_country = ?)", alias, alias))
		args = append(args, strings.ToUpper(f.Country), strings.ToUpper(f.Country))
	}
	if f.MinBytes > 0 {
		clauses = append(clauses, fmt.Sprintf("%sbytes >= ?", alias))
		args = append(args, f.MinBytes)
	}
	if f.Search != "" {
		needle := "%" + strings.ToLower(f.Search) + "%"
		clauses = append(clauses,
			fmt.Sprintf("(lower(%ssrc_ip) LIKE ? OR lower(%sdst_ip) LIKE ? OR lower(%ssrc_hostname) LIKE ? OR lower(%sdst_hostname) LIKE ? OR lower(%ssrc_service) LIKE ? OR lower(%sdst_service) LIKE ?)", alias, alias, alias, alias, alias, alias),
		)
		args = append(args, needle, needle, needle, needle, needle, needle)
	}
	return strings.Join(clauses, " AND "), args
}

func safeSort(input string, allowed []string, fallback string) string {
	in := strings.TrimSpace(input)
	for _, v := range allowed {
		if in == v {
			return in
		}
	}
	return fallback
}

func safeSortDir(input string) string {
	v := strings.ToUpper(strings.TrimSpace(input))
	if v == "ASC" {
		return "ASC"
	}
	return "DESC"
}
