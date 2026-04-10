package storage

import (
	"testing"
	"time"

	"github.com/flowscope/flowscope/internal/model"
)

func TestBuildFlowWhere(t *testing.T) {
	f := model.QueryFilter{
		From:       time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC),
		To:         time.Date(2026, 4, 10, 11, 0, 0, 0, time.UTC),
		SrcIP:      "10.10.1.10",
		Protocol:   "tcp",
		DstPort:    443,
		MinBytes:   1000,
		Country:    "US",
		ExporterID: "exp-a",
	}
	where, args := buildFlowWhere(f, "r")
	if where == "" {
		t.Fatal("where empty")
	}
	if len(args) < 8 {
		t.Fatalf("expected args >= 8, got %d", len(args))
	}
}

func TestSafeSort(t *testing.T) {
	if got := safeSort("bytes", []string{"bytes", "packets"}, "packets"); got != "bytes" {
		t.Fatalf("got %s", got)
	}
	if got := safeSort("invalid", []string{"bytes", "packets"}, "packets"); got != "packets" {
		t.Fatalf("got %s", got)
	}
}
