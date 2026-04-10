package normalize

import (
	"testing"
	"time"

	"github.com/flowscope/flowscope/internal/model"
)

func TestApplySetsDerivedFields(t *testing.T) {
	rec := model.FlowRecord{
		TimestampStart: time.Date(2026, 4, 10, 9, 12, 33, 0, time.UTC),
		ExporterID:     "exp-a",
		SrcIP:          "10.10.1.10",
		DstIP:          "198.51.100.4",
		IPProtocol:     6,
		SrcPort:        12345,
		DstPort:        443,
	}
	Apply(&rec)
	if rec.FlowID == "" {
		t.Fatalf("expected flow id")
	}
	if rec.SrcSubnet == "" || rec.DstSubnet == "" {
		t.Fatalf("expected subnets")
	}
	if rec.FlowKeyHash == "" {
		t.Fatalf("expected hash")
	}
	if rec.L4ProtocolName != "TCP" {
		t.Fatalf("expected TCP, got %s", rec.L4ProtocolName)
	}
	if rec.FlowDirection != "egress" {
		t.Fatalf("expected egress, got %s", rec.FlowDirection)
	}
}
