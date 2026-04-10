package model

import "time"

type FlowRecord struct {
	FlowID          string    `json:"flow_id"`
	TimestampStart  time.Time `json:"timestamp_start"`
	TimestampEnd    time.Time `json:"timestamp_end"`
	ExporterID      string    `json:"exporter_id"`
	ExporterIP      string    `json:"exporter_ip"`
	ExporterName    string    `json:"exporter_name,omitempty"`
	ObservationID   uint32    `json:"observation_domain,omitempty"`
	SrcIP           string    `json:"src_ip"`
	DstIP           string    `json:"dst_ip"`
	SrcPort         uint16    `json:"src_port"`
	DstPort         uint16    `json:"dst_port"`
	IPProtocol      uint8     `json:"ip_protocol"`
	L4ProtocolName  string    `json:"l4_protocol_name"`
	Bytes           uint64    `json:"bytes"`
	Packets         uint64    `json:"packets"`
	InputInterface  uint32    `json:"input_interface"`
	OutputInterface uint32    `json:"output_interface"`
	InputIfAlias    string    `json:"input_interface_alias,omitempty"`
	OutputIfAlias   string    `json:"output_interface_alias,omitempty"`
	SrcASN          uint32    `json:"src_asn"`
	DstASN          uint32    `json:"dst_asn"`
	SrcCountry      string    `json:"src_country"`
	DstCountry      string    `json:"dst_country"`
	SrcHostname     string    `json:"src_hostname,omitempty"`
	DstHostname     string    `json:"dst_hostname,omitempty"`
	SrcService      string    `json:"src_service,omitempty"`
	DstService      string    `json:"dst_service,omitempty"`
	SrcEnvironment  string    `json:"src_environment,omitempty"`
	DstEnvironment  string    `json:"dst_environment,omitempty"`
	SrcOwnerTeam    string    `json:"src_owner_team,omitempty"`
	DstOwnerTeam    string    `json:"dst_owner_team,omitempty"`
	SrcMAC          string    `json:"src_mac,omitempty"`
	DstMAC          string    `json:"dst_mac,omitempty"`
	VLANID          uint16    `json:"vlan_id,omitempty"`
	TCPFlags        uint16    `json:"tcp_flags,omitempty"`
	FlowDirection   string    `json:"flow_direction,omitempty"`
	SamplerRate     uint32    `json:"sampler_rate,omitempty"`
	SrcSubnet       string    `json:"src_subnet"`
	DstSubnet       string    `json:"dst_subnet"`
	SrcIsPrivate    bool      `json:"src_is_private"`
	DstIsPrivate    bool      `json:"dst_is_private"`
	FlowKeyHash     string    `json:"flow_key_hash"`
	MinuteBucket    time.Time `json:"minute_bucket"`
	HourBucket      time.Time `json:"hour_bucket"`
	SourceType      string    `json:"source_type"`
}

type Exporter struct {
	ExporterID     string    `json:"exporter_id"`
	ExporterIP     string    `json:"exporter_ip"`
	ExporterName   string    `json:"exporter_name"`
	ObservationID  uint32    `json:"observation_domain"`
	FirstSeen      time.Time `json:"first_seen"`
	LastSeen       time.Time `json:"last_seen"`
	FlowsObserved  uint64    `json:"flows_observed"`
	LastSourceType string    `json:"last_source_type"`
}

type InventoryAsset struct {
	AssetID      string   `json:"asset_id" yaml:"asset_id"`
	CIDR         string   `json:"cidr" yaml:"cidr"`
	Hostname     string   `json:"hostname" yaml:"hostname"`
	Service      string   `json:"service" yaml:"service"`
	Environment  string   `json:"environment" yaml:"environment"`
	OwnerTeam    string   `json:"owner_team" yaml:"owner_team"`
	InterfaceMap []string `json:"interface_aliases,omitempty" yaml:"interface_aliases,omitempty"`
}

type QueryFilter struct {
	From            time.Time
	To              time.Time
	Page            int
	PageSize        int
	SortBy          string
	SortDir         string
	SrcIP           string
	DstIP           string
	Subnet          string
	ExporterID      string
	Protocol        string
	SrcPort         int
	DstPort         int
	InputInterface  int
	OutputInterface int
	ASN             int
	Country         string
	MinBytes        uint64
	MinFlows        uint64
	Search          string
	Mode            string
	GroupBy         string
	NodeID          string
	EdgeID          string
}

type PageResult[T any] struct {
	Data       []T  `json:"data"`
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	Total      uint64 `json:"total"`
	HasNext    bool `json:"has_next"`
	SortBy     string `json:"sort_by"`
	SortDir    string `json:"sort_dir"`
}

type TopItem struct {
	Key     string `json:"key"`
	Bytes   uint64 `json:"bytes"`
	Packets uint64 `json:"packets"`
	Flows   uint64 `json:"flows"`
}

type TimePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Bytes     uint64    `json:"bytes"`
	Packets   uint64    `json:"packets"`
	Flows     uint64    `json:"flows"`
}

type GraphNode struct {
	ID        string            `json:"id"`
	Label     string            `json:"label"`
	Type      string            `json:"type"`
	BytesIn   uint64            `json:"bytes_in"`
	BytesOut  uint64            `json:"bytes_out"`
	PacketsIn uint64            `json:"packets_in"`
	PacketsOut uint64           `json:"packets_out"`
	FlowCount uint64            `json:"flow_count"`
	LastSeen  time.Time         `json:"last_seen"`
	Tags      map[string]string `json:"tags"`
	Private   bool              `json:"private"`
	Internal  bool              `json:"internal"`
	Collapsed bool              `json:"collapsed"`
	Children  []string          `json:"children,omitempty"`
}

type GraphEdge struct {
	ID          string    `json:"id"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Bytes       uint64    `json:"bytes"`
	Packets     uint64    `json:"packets"`
	Flows       uint64    `json:"flows"`
	Protocols   []string  `json:"protocols"`
	TopDstPorts []uint16  `json:"top_destination_ports"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
	Exporters   uint64    `json:"exporter_count"`
	Direction   string    `json:"direction"`
}

type GraphResponse struct {
	Nodes       []GraphNode          `json:"nodes"`
	Edges       []GraphEdge          `json:"edges"`
	Mode        string               `json:"mode"`
	GroupBy     string               `json:"group_by"`
	Suspicious  []SuspiciousInteraction `json:"suspicious_interactions"`
}

type SuspiciousInteraction struct {
	Kind        string `json:"kind"`
	Description string `json:"description"`
	NodeID      string `json:"node_id,omitempty"`
	EdgeID      string `json:"edge_id,omitempty"`
	Severity    string `json:"severity"`
}

type NodeDetails struct {
	Node        GraphNode   `json:"node"`
	Inbound     TopItem     `json:"inbound_totals"`
	Outbound    TopItem     `json:"outbound_totals"`
	TopPeers    []TopItem   `json:"top_peers"`
	TopPorts    []TopItem   `json:"top_ports"`
	TopProtocols []TopItem  `json:"top_protocols"`
	Series      []TimePoint `json:"timeseries"`
}

type EdgeDetails struct {
	Edge         GraphEdge   `json:"edge"`
	Directional  []TopItem   `json:"directional_stats"`
	Series       []TimePoint `json:"timeseries"`
	Protocols    []TopItem   `json:"protocols"`
	PortDist     []TopItem   `json:"port_distribution"`
	Exporters    []TopItem   `json:"exporters"`
	FlowPage     PageResult[FlowRecord] `json:"flows"`
}

type SankeyNode struct {
	Name string `json:"name"`
}

type SankeyLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  uint64 `json:"value"`
	Flows  uint64 `json:"flows"`
}

type SankeyResponse struct {
	Nodes []SankeyNode `json:"nodes"`
	Links []SankeyLink `json:"links"`
}
