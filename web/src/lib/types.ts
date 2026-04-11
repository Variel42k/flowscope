export type FlowRecord = {
  flow_id: string
  timestamp_start: string
  timestamp_end: string
  exporter_id: string
  exporter_ip: string
  exporter_name?: string
  src_ip: string
  dst_ip: string
  src_port: number
  dst_port: number
  ip_protocol: number
  l4_protocol_name: string
  bytes: number
  packets: number
  input_interface: number
  output_interface: number
  input_interface_alias?: string
  output_interface_alias?: string
  src_asn: number
  dst_asn: number
  src_country: string
  dst_country: string
  src_hostname?: string
  dst_hostname?: string
  src_service?: string
  dst_service?: string
  src_environment?: string
  dst_environment?: string
  src_owner_team?: string
  dst_owner_team?: string
  flow_direction?: string
  sampler_rate?: number
  src_subnet: string
  dst_subnet: string
  src_is_private: boolean
  dst_is_private: boolean
  source_type: string
}

export type PageResult<T> = {
  data: T[]
  page: number
  page_size: number
  total: number
  has_next: boolean
  sort_by: string
  sort_dir: string
}

export type TopItem = {
  key: string
  bytes: number
  packets: number
  flows: number
}

export type SankeyResponse = {
  nodes: { name: string }[]
  links: { source: string; target: string; value: number; flows: number }[]
}

export type GraphNode = {
  id: string
  label: string
  type: string
  bytes_in: number
  bytes_out: number
  packets_in: number
  packets_out: number
  flow_count: number
  last_seen: string
  tags: Record<string, string>
  private: boolean
  internal: boolean
  collapsed?: boolean
  children?: string[]
}

export type GraphEdge = {
  id: string
  source: string
  destination: string
  bytes: number
  packets: number
  flows: number
  protocols: string[]
  top_destination_ports: number[]
  first_seen: string
  last_seen: string
  exporter_count: number
  direction: string
}

export type SuspiciousInteraction = {
  kind: string
  description: string
  node_id?: string
  edge_id?: string
  severity: string
}

export type GraphResponse = {
  nodes: GraphNode[]
  edges: GraphEdge[]
  mode: string
  group_by: string
  suspicious_interactions: SuspiciousInteraction[]
}

export type NodeDetails = {
  node: GraphNode
  inbound_totals: TopItem
  outbound_totals: TopItem
  top_peers: TopItem[]
  top_ports: TopItem[]
  top_protocols: TopItem[]
  timeseries: { timestamp: string; bytes: number; packets: number; flows: number }[]
}

export type EdgeDetails = {
  edge: GraphEdge
  directional_stats: TopItem[]
  timeseries: { timestamp: string; bytes: number; packets: number; flows: number }[]
  protocols: TopItem[]
  port_distribution: TopItem[]
  exporters: TopItem[]
  flows: PageResult<FlowRecord>
}

export type AuthSession = {
  token: string
  user: string
  role: string
  auth_mode?: 'local' | 'oidc'
}

export type AlertRule = {
  rule_id: string
  name: string
  rule_type: 'new_edge' | 'fanout_external' | 'high_byte_edge' | 'port_outlier'
  enabled: boolean
  threshold_value: number
  window_minutes: number
  severity: 'low' | 'medium' | 'high'
  created_by: string
  created_at: string
  updated_at: string
}

export type AlertEvent = {
  event_id: string
  event_key: string
  rule_id: string
  rule_name: string
  rule_type: string
  severity: string
  detected_at: string
  window_from: string
  window_to: string
  node_id?: string
  edge_id?: string
  description: string
  bytes: number
  flows: number
  metadata?: Record<string, unknown>
}

export type SavedView = {
  view_id: string
  name: string
  description: string
  scope: 'overview' | 'flows' | 'sankey' | 'map' | 'global'
  owner_user: string
  is_shared: boolean
  filters: Record<string, string>
  created_at: string
  updated_at: string
}
