CREATE DATABASE IF NOT EXISTS flowscope;

CREATE TABLE IF NOT EXISTS flowscope.raw_flow_events (
    flow_id String,
    timestamp_start DateTime64(3, 'UTC'),
    timestamp_end DateTime64(3, 'UTC'),
    exporter_id String,
    exporter_ip String,
    exporter_name String,
    observation_id UInt32,
    src_ip String,
    dst_ip String,
    src_port UInt16,
    dst_port UInt16,
    ip_protocol UInt8,
    l4_protocol_name LowCardinality(String),
    bytes UInt64,
    packets UInt64,
    input_interface UInt32,
    output_interface UInt32,
    input_interface_alias String,
    output_interface_alias String,
    src_asn UInt32,
    dst_asn UInt32,
    src_country String,
    dst_country String,
    src_hostname String,
    dst_hostname String,
    src_service String,
    dst_service String,
    src_environment String,
    dst_environment String,
    src_owner_team String,
    dst_owner_team String,
    src_mac String,
    dst_mac String,
    vlan_id UInt16,
    tcp_flags UInt16,
    flow_direction LowCardinality(String),
    sampler_rate UInt32,
    src_subnet String,
    dst_subnet String,
    src_is_private Bool,
    dst_is_private Bool,
    flow_key_hash String,
    minute_bucket DateTime('UTC'),
    hour_bucket DateTime('UTC'),
    source_type LowCardinality(String)
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(minute_bucket)
ORDER BY (minute_bucket, exporter_id, src_ip, dst_ip, dst_port, ip_protocol)
TTL minute_bucket + INTERVAL 30 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS flowscope.flow_1m_rollup (
    minute_bucket DateTime('UTC'),
    exporter_id String,
    src_subnet String,
    dst_subnet String,
    protocol LowCardinality(String),
    dst_port UInt16,
    input_interface UInt32,
    output_interface UInt32,
    bytes UInt64,
    packets UInt64,
    flows UInt64
)
ENGINE = SummingMergeTree
PARTITION BY toYYYYMMDD(minute_bucket)
ORDER BY (minute_bucket, exporter_id, src_subnet, dst_subnet, protocol, dst_port)
TTL minute_bucket + INTERVAL 30 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS flowscope.flow_1h_rollup (
    hour_bucket DateTime('UTC'),
    exporter_id String,
    src_subnet String,
    dst_subnet String,
    protocol LowCardinality(String),
    dst_port UInt16,
    bytes UInt64,
    packets UInt64,
    flows UInt64
)
ENGINE = SummingMergeTree
PARTITION BY toYYYYMMDD(hour_bucket)
ORDER BY (hour_bucket, exporter_id, src_subnet, dst_subnet, protocol, dst_port)
TTL hour_bucket + INTERVAL 180 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS flowscope.nodes_dimension (
    node_id String,
    node_type LowCardinality(String),
    bytes_in UInt64,
    bytes_out UInt64,
    packets_in UInt64,
    packets_out UInt64,
    flows UInt64,
    last_seen DateTime('UTC'),
    tags_json String,
    updated_at DateTime('UTC')
)
ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY toYYYYMMDD(updated_at)
ORDER BY (node_id)
TTL updated_at + INTERVAL 30 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS flowscope.edges_1m (
    minute_bucket DateTime('UTC'),
    exporter_id String,
    src_node_id String,
    dst_node_id String,
    src_subnet String,
    dst_subnet String,
    src_service String,
    dst_service String,
    src_environment String,
    dst_environment String,
    src_asn UInt32,
    dst_asn UInt32,
    src_country String,
    dst_country String,
    src_private Bool,
    dst_private Bool,
    protocol LowCardinality(String),
    dst_port UInt16,
    bytes UInt64,
    packets UInt64,
    flows UInt64,
    first_seen DateTime64(3, 'UTC'),
    last_seen DateTime64(3, 'UTC')
)
ENGINE = SummingMergeTree
PARTITION BY toYYYYMMDD(minute_bucket)
ORDER BY (minute_bucket, src_node_id, dst_node_id, protocol, dst_port, exporter_id)
TTL minute_bucket + INTERVAL 30 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS flowscope.edges_1h (
    hour_bucket DateTime('UTC'),
    exporter_id String,
    src_node_id String,
    dst_node_id String,
    src_subnet String,
    dst_subnet String,
    src_service String,
    dst_service String,
    src_environment String,
    dst_environment String,
    src_asn UInt32,
    dst_asn UInt32,
    src_country String,
    dst_country String,
    src_private Bool,
    dst_private Bool,
    protocol LowCardinality(String),
    dst_port UInt16,
    bytes UInt64,
    packets UInt64,
    flows UInt64,
    first_seen DateTime64(3, 'UTC'),
    last_seen DateTime64(3, 'UTC')
)
ENGINE = SummingMergeTree
PARTITION BY toYYYYMMDD(hour_bucket)
ORDER BY (hour_bucket, src_node_id, dst_node_id, protocol, dst_port, exporter_id)
TTL hour_bucket + INTERVAL 180 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS flowscope.exporters (
    exporter_id String,
    exporter_ip String,
    exporter_name String,
    observation_id UInt32,
    first_seen DateTime('UTC'),
    last_seen DateTime('UTC'),
    flows_observed UInt64,
    last_source_type LowCardinality(String)
)
ENGINE = ReplacingMergeTree(last_seen)
PARTITION BY toYYYYMMDD(last_seen)
ORDER BY (exporter_id, last_seen)
TTL last_seen + INTERVAL 365 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS flowscope.inventory_assets (
    asset_id String,
    cidr String,
    hostname String,
    service String,
    environment String,
    owner_team String,
    interface_aliases String,
    updated_at DateTime64(3, 'UTC')
)
ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY toYYYYMMDD(updated_at)
ORDER BY (cidr, asset_id)
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS flowscope.alert_rules (
    rule_id String,
    name String,
    rule_type LowCardinality(String),
    enabled Bool,
    threshold_value UInt64,
    window_minutes UInt32,
    severity LowCardinality(String),
    created_by String,
    created_at DateTime64(3, 'UTC'),
    updated_at DateTime64(3, 'UTC'),
    deleted Bool
)
ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY toYYYYMMDD(updated_at)
ORDER BY (rule_id)
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS flowscope.alert_events (
    event_id String,
    event_key String,
    rule_id String,
    rule_name String,
    rule_type LowCardinality(String),
    severity LowCardinality(String),
    detected_at DateTime64(3, 'UTC'),
    window_from DateTime('UTC'),
    window_to DateTime('UTC'),
    node_id String,
    edge_id String,
    description String,
    bytes UInt64,
    flows UInt64,
    metadata_json String
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(detected_at)
ORDER BY (detected_at, severity, rule_id, event_id)
TTL detected_at + INTERVAL 180 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS flowscope.saved_views (
    view_id String,
    name String,
    description String,
    scope LowCardinality(String),
    owner_user String,
    is_shared Bool,
    filters_json String,
    created_at DateTime64(3, 'UTC'),
    updated_at DateTime64(3, 'UTC'),
    deleted Bool
)
ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY toYYYYMMDD(updated_at)
ORDER BY (scope, owner_user, view_id)
TTL updated_at + INTERVAL 365 DAY
SETTINGS index_granularity = 8192;
