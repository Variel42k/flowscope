package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr             string
	ClickHouseDSN        string
	JWTSecret            string
	AdminUser            string
	AdminPassword        string
	AdminRole            string
	ViewerRole           string
	AuthDisabled         bool
	LogLevel             string
	RetentionDays        int
	OIDCEnabled          bool
	OIDCIssuerURL        string
	OIDCClientID         string
	OIDCClientSecret     string
	OIDCRedirectURL      string
	OIDCSuccessRedirect  string
	OIDCScopes           string
	OIDCRoleClaim        string
	OIDCAdminUsers       []string
	OIDCAllowedDomains   []string
	InventoryPath        string
	InventoryFormat      string
	GeoIPCityPath        string
	GeoIPASNPath         string
	RDNSCacheTTL         time.Duration
	EnableGeoIP          bool
	EnableRDNS           bool
	EnableInventory      bool
	EnableInterfaceAlias bool
	InterfaceAliasPath   string
	CollectorNetFlowV5   string
	CollectorNetFlowV9   string
	CollectorIPFIX       string
	CollectorSFlow       string
	QueryDefaultWindow   time.Duration
	RollupInterval       time.Duration
}

func Load() Config {
	cfg := Config{
		HTTPAddr:             env("FLOWSCOPE_HTTP_ADDR", ":8088"),
		ClickHouseDSN:        env("FLOWSCOPE_CLICKHOUSE_DSN", "clickhouse://default:flowscope@clickhouse:9000/flowscope"),
		JWTSecret:            env("FLOWSCOPE_JWT_SECRET", "flowscope-dev-secret"),
		AdminUser:            env("FLOWSCOPE_ADMIN_USER", "admin"),
		AdminPassword:        env("FLOWSCOPE_ADMIN_PASSWORD", "admin123"),
		AdminRole:            env("FLOWSCOPE_ADMIN_ROLE", "admin"),
		ViewerRole:           env("FLOWSCOPE_VIEWER_ROLE", "viewer"),
		AuthDisabled:         envBool("FLOWSCOPE_AUTH_DISABLED", false),
		LogLevel:             env("FLOWSCOPE_LOG_LEVEL", "info"),
		RetentionDays:        envInt("FLOWSCOPE_RETENTION_DAYS", 30),
		OIDCEnabled:          envBool("FLOWSCOPE_OIDC_ENABLED", false),
		OIDCIssuerURL:        env("FLOWSCOPE_OIDC_ISSUER_URL", ""),
		OIDCClientID:         env("FLOWSCOPE_OIDC_CLIENT_ID", ""),
		OIDCClientSecret:     env("FLOWSCOPE_OIDC_CLIENT_SECRET", ""),
		OIDCRedirectURL:      env("FLOWSCOPE_OIDC_REDIRECT_URL", ""),
		OIDCSuccessRedirect:  env("FLOWSCOPE_OIDC_SUCCESS_REDIRECT", "http://localhost:5173/oidc/callback"),
		OIDCScopes:           env("FLOWSCOPE_OIDC_SCOPES", "openid,profile,email"),
		OIDCRoleClaim:        env("FLOWSCOPE_OIDC_ROLE_CLAIM", "role"),
		OIDCAdminUsers:       envCSV("FLOWSCOPE_OIDC_ADMIN_USERS"),
		OIDCAllowedDomains:   envCSV("FLOWSCOPE_OIDC_ALLOWED_DOMAINS"),
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

func envCSV(k string) []string {
	v := strings.TrimSpace(env(k, ""))
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			out = append(out, p)
		}
	}
	return out
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
