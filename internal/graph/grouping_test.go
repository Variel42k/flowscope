package graph

import (
	"testing"
	"time"

	"github.com/flowscope/flowscope/internal/model"
)

func TestCollapseByEnvironment(t *testing.T) {
	in := model.GraphResponse{
		Nodes: []model.GraphNode{
			{ID: "10.10.1.10", Label: "10.10.1.10", Tags: map[string]string{"environment": "prod"}, FlowCount: 5, LastSeen: time.Now()},
			{ID: "10.10.2.11", Label: "10.10.2.11", Tags: map[string]string{"environment": "prod"}, FlowCount: 3, LastSeen: time.Now()},
			{ID: "198.51.100.8", Label: "198.51.100.8", Tags: map[string]string{"environment": "external"}, FlowCount: 10, LastSeen: time.Now()},
		},
		Edges: []model.GraphEdge{{ID: "a", Source: "10.10.1.10", Destination: "198.51.100.8", Bytes: 100, Flows: 1}},
	}
	out := Collapse(in, "environment")
	if len(out.Nodes) != 2 {
		t.Fatalf("expected 2 grouped nodes, got %d", len(out.Nodes))
	}
	if len(out.Edges) != 1 {
		t.Fatalf("expected 1 grouped edge, got %d", len(out.Edges))
	}
}
