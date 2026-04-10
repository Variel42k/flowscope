package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/flowscope/flowscope/internal/model"
)

func parseFilter(r *http.Request) model.QueryFilter {
	q := r.URL.Query()
	f := model.QueryFilter{
		From:            parseTime(q.Get("from")),
		To:              parseTime(q.Get("to")),
		Page:            parseInt(q.Get("page"), 1),
		PageSize:        parseInt(q.Get("page_size"), 50),
		SortBy:          q.Get("sort_by"),
		SortDir:         q.Get("sort_dir"),
		SrcIP:           q.Get("src_ip"),
		DstIP:           q.Get("dst_ip"),
		Subnet:          q.Get("subnet"),
		ExporterID:      q.Get("exporter"),
		Protocol:        q.Get("protocol"),
		SrcPort:         parseInt(q.Get("src_port"), 0),
		DstPort:         parseInt(q.Get("dst_port"), 0),
		InputInterface:  parseInt(q.Get("input_interface"), 0),
		OutputInterface: parseInt(q.Get("output_interface"), 0),
		ASN:             parseInt(q.Get("asn"), 0),
		Country:         q.Get("country"),
		MinBytes:        uint64(parseInt(q.Get("min_bytes"), 0)),
		MinFlows:        uint64(parseInt(q.Get("min_flows"), 0)),
		Search:          q.Get("search"),
		Mode:            q.Get("mode"),
		GroupBy:         q.Get("group_by"),
		NodeID:          q.Get("node_id"),
		EdgeID:          q.Get("edge_id"),
	}
	if f.From.IsZero() {
		f.From = time.Now().UTC().Add(-15 * time.Minute)
	}
	if f.To.IsZero() {
		f.To = time.Now().UTC()
	}
	if f.Mode == "" {
		f.Mode = "host_to_host"
	}
	if f.GroupBy == "" {
		f.GroupBy = "none"
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 50
	}
	return f
}

func parseInt(v string, def int) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func parseTime(v string) time.Time {
	if v == "" {
		return time.Time{}
	}
	if t, err := time.Parse(time.RFC3339, v); err == nil {
		return t.UTC()
	}
	if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
		return t.UTC()
	}
	return time.Time{}
}
