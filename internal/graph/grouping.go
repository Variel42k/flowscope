package graph

import (
	"strings"
	"time"

	"github.com/flowscope/flowscope/internal/model"
	"github.com/flowscope/flowscope/internal/util"
)

func Collapse(input model.GraphResponse, groupBy string) model.GraphResponse {
	if groupBy == "" || groupBy == "none" {
		return input
	}
	mapNode := map[string]string{}
	groups := map[string]*model.GraphNode{}
	for _, n := range input.Nodes {
		key := groupKey(n, groupBy)
		mapNode[n.ID] = key
		g, ok := groups[key]
		if !ok {
			g = &model.GraphNode{ID: key, Label: key, Type: groupBy, Tags: map[string]string{}, Collapsed: true}
			groups[key] = g
		}
		g.Children = append(g.Children, n.ID)
		g.BytesIn += n.BytesIn
		g.BytesOut += n.BytesOut
		g.PacketsIn += n.PacketsIn
		g.PacketsOut += n.PacketsOut
		g.FlowCount += n.FlowCount
		if n.LastSeen.After(g.LastSeen) {
			g.LastSeen = n.LastSeen
		}
		for k, v := range n.Tags {
			if g.Tags[k] == "" {
				g.Tags[k] = v
			}
		}
	}
	edges := map[string]*model.GraphEdge{}
	for _, e := range input.Edges {
		src := mapNode[e.Source]
		dst := mapNode[e.Destination]
		if src == "" || dst == "" {
			continue
		}
		id := src + "->" + dst
		ge, ok := edges[id]
		if !ok {
			clone := e
			clone.ID = id
			clone.Source = src
			clone.Destination = dst
			clone.Protocols = append([]string{}, e.Protocols...)
			clone.TopDstPorts = append([]uint16{}, e.TopDstPorts...)
			edges[id] = &clone
			continue
		}
		ge.Bytes += e.Bytes
		ge.Packets += e.Packets
		ge.Flows += e.Flows
		if e.FirstSeen.Before(ge.FirstSeen) || ge.FirstSeen.IsZero() {
			ge.FirstSeen = e.FirstSeen
		}
		if e.LastSeen.After(ge.LastSeen) {
			ge.LastSeen = e.LastSeen
		}
		ge.Exporters += e.Exporters
		ge.Protocols = mergeStrings(ge.Protocols, e.Protocols)
		ge.TopDstPorts = mergePorts(ge.TopDstPorts, e.TopDstPorts)
	}
	out := model.GraphResponse{Mode: input.Mode, GroupBy: groupBy, Suspicious: input.Suspicious}
	for _, n := range groups {
		out.Nodes = append(out.Nodes, *n)
	}
	for _, e := range edges {
		out.Edges = append(out.Edges, *e)
	}
	return out
}

func groupKey(n model.GraphNode, groupBy string) string {
	switch groupBy {
	case "subnet":
		if strings.Contains(n.ID, "/") {
			return n.ID
		}
		return util.SubnetForIP(n.ID)
	case "service":
		if v := n.Tags["service"]; v != "" {
			return "service:" + v
		}
		return "service:unknown"
	case "environment":
		if v := n.Tags["environment"]; v != "" {
			return "env:" + v
		}
		return "env:unknown"
	case "asn":
		if v := n.Tags["asn"]; v != "" {
			return "asn:" + v
		}
		return "asn:unknown"
	default:
		return n.ID
	}
}

func mergeStrings(a, b []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(a)+len(b))
	for _, x := range a {
		if _, ok := seen[x]; ok {
			continue
		}
		seen[x] = struct{}{}
		out = append(out, x)
	}
	for _, x := range b {
		if _, ok := seen[x]; ok {
			continue
		}
		seen[x] = struct{}{}
		out = append(out, x)
	}
	return out
}

func mergePorts(a, b []uint16) []uint16 {
	seen := map[uint16]struct{}{}
	out := make([]uint16, 0, len(a)+len(b))
	for _, x := range a {
		if _, ok := seen[x]; ok {
			continue
		}
		seen[x] = struct{}{}
		out = append(out, x)
	}
	for _, x := range b {
		if _, ok := seen[x]; ok {
			continue
		}
		seen[x] = struct{}{}
		out = append(out, x)
	}
	if len(out) > 10 {
		out = out[:10]
	}
	return out
}

func Expand(input model.GraphResponse, groupNodeID string) model.GraphResponse {
	_ = time.Now()
	// Expansion is handled client-side using node.children metadata and cached original graph.
	return input
}
