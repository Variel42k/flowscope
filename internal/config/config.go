package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr              string
	ClickHouseDSN         string
	JWTSecret             string
	AdminUser             string
	AdminPassword         string
	AuthDisabled          bool
	LogLevel              string
	RetentionDays         int
	InventoryPath         string
	InventoryFormat       string
	GeoIPCityPath         string
	GeoIPASNPath          string
	RDNSCacheTTL          time.Duration
	EnableGeoIP           bool
	EnableRDNS            bool
	EnableInventory       bool
	EnableInterfaceAlias  bool
	InterfaceAliasPath    string
	CollectorNetFlowV5    string
	CollectorNetFlowV9    string
	CollectorIPFIX        string
	CollectorSFlow        string
	QueryDefaultWindow    time.Duration
	RollupInterval        time.Duration
}

func Load() Config {
	cfg := Config{
		HTTPAddr:             env("FLOWSCOPE_HTTP_ADDR", ":8088"),
		ClickHouseDSN:        env("FLOWSCOPE_CLICKHOUSE_DSN", "clickhouse://default:flowscope@clickhouse:9000/flowscope"),
		JWTSecret:            env("FLOWSCOPE_JWT_SECRET", "flowscope-dev-secret"),
		AdminUser:            env("FLOWSCOPE_ADMIN_USER", "admin"),
		AdminPassword:        env("FLOWSCOPE_ADMIN_PASSWORD", "admin123"),
		AuthDisabled:         envBool("FLOWSCOPE_AUTH_DISABLED", false),
		LogLevel:             env("FLOWSCOPE_LOG_LEVEL", "info"),
		RetentionDays:        envInt("FLOWSCOPE_RETENTION_DAYS", 30),
		InventoryPath:        env("FLOWSCOPE_INVENTORY_PATH", "/data/inventory.sample.yaml"),
		InventoryFormat:      env("FLOWSCOPE_INVENTORY_FORMAT", "yaml"),
		GeoIPCityPath:        env("FLOWSCOPE_GEOIP_CITY_PATH", ""),
		GeoIPASNPath:         env("FLOWSCOPE_GEOIP_ASN_PATH", ""),
		RDNSCacheTTL:         envDuration("FLOWSCOPE_RDNS_CACHE_TTL", 30*time.Minute),
		EnableGeoIP:          envBool("FLOWSCOPE_ENABLE_GEOIP", false),
		EnableRDNS:           envBool("FLOWSCOPE_ENABLE_RDNS", false),
		EnableInventory:      envBool("FLOWSCOPE_ENABLE_INVENTORY", true),
		EnableInterfaceAlias: envBool("FLOWSCOPE_ENABLE_INTERFACE_ALIAS", true),
		InterfaceAliasPath:   env("FLOWSCOPE_INTERFACE_ALIAS_PATH", "/data/interfaces.sample.yaml"),
		CollectorNetFlowV5:   env("FLOWSCOPE_LISTEN_NETFLOW_V5", ":2055"),
		CollectorNetFlowV9:   env("FLOWSCOPE_LISTEN_NETFLOW_V9", ":2056"),
		CollectorIPFIX:       env("FLOWSCOPE_LISTEN_IPFIX", ":4739"),
		CollectorSFlow:       env("FLOWSCOPE_LISTEN_SFLOW", ":6343"),
		QueryDefaultWindow:   envDuration("FLOWSCOPE_DEFAULT_QUERY_WINDOW", 15*time.Minute),
		RollupInterval:       envDuration("FLOWSCOPE_ROLLUP_INTERVAL", time.Minute),
	}
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		log.Fatal("FLOWSCOPE_JWT_SECRET must not be empty")
	}
	return cfg
}

func env(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
}

func envInt(k string, def int) int {
	if v, ok := os.LookupEnv(k); ok {
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
	}
	return def
}

func envBool(k string, def bool) bool {
	if v, ok := os.LookupEnv(k); ok {
		v = strings.ToLower(strings.TrimSpace(v))
		return v == "1" || v == "true" || v == "yes" || v == "on"
	}
	return def
}

func envDuration(k string, def time.Duration) time.Duration {
	if v, ok := os.LookupEnv(k); ok {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return def
}
