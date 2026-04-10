package util

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

func ProtocolName(proto uint8) string {
	switch proto {
	case 1:
		return "ICMP"
	case 6:
		return "TCP"
	case 17:
		return "UDP"
	case 47:
		return "GRE"
	case 50:
		return "ESP"
	case 58:
		return "ICMPv6"
	default:
		return fmt.Sprintf("IP-%d", proto)
	}
}

func SubnetForIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "unknown"
	}
	if v4 := ip.To4(); v4 != nil {
		mask := net.CIDRMask(24, 32)
		netIP := v4.Mask(mask)
		return fmt.Sprintf("%s/24", netIP.String())
	}
	mask := net.CIDRMask(64, 128)
	netIP := ip.Mask(mask)
	return fmt.Sprintf("%s/64", netIP.String())
}

func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"100.64.0.0/10",
		"fc00::/7",
		"fe80::/10",
	}
	for _, cidr := range privateCIDRs {
		_, n, err := net.ParseCIDR(cidr)
		if err == nil && n.Contains(ip) {
			return true
		}
	}
	return false
}

func FlowHash(parts ...any) string {
	builder := strings.Builder{}
	for _, p := range parts {
		builder.WriteString(fmt.Sprint(p))
		builder.WriteByte('|')
	}
	h := sha1.Sum([]byte(builder.String()))
	return hex.EncodeToString(h[:])
}

func MinuteBucket(t time.Time) time.Time {
	return t.UTC().Truncate(time.Minute)
}

func HourBucket(t time.Time) time.Time {
	return t.UTC().Truncate(time.Hour)
}
