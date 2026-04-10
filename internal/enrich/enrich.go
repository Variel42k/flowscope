package enrich

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/flowscope/flowscope/internal/inventory"
	"github.com/flowscope/flowscope/internal/model"
	"github.com/oschwald/geoip2-golang"
)

type Enricher interface {
	Enrich(ctx context.Context, record *model.FlowRecord)
}

type Composite struct {
	items []Enricher
}

func NewComposite(items ...Enricher) *Composite {
	return &Composite{items: items}
}

func (c *Composite) Enrich(ctx context.Context, record *model.FlowRecord) {
	for _, item := range c.items {
		if item != nil {
			item.Enrich(ctx, record)
		}
	}
}

type GeoIPEnricher struct {
	city *geoip2.Reader
	asn  *geoip2.Reader
}

func NewGeoIPEnricher(cityPath string, asnPath string) (*GeoIPEnricher, error) {
	g := &GeoIPEnricher{}
	if strings.TrimSpace(cityPath) != "" {
		reader, err := geoip2.Open(cityPath)
		if err != nil {
			return nil, err
		}
		g.city = reader
	}
	if strings.TrimSpace(asnPath) != "" {
		reader, err := geoip2.Open(asnPath)
		if err != nil {
			return nil, err
		}
		g.asn = reader
	}
	return g, nil
}

func (g *GeoIPEnricher) Enrich(_ context.Context, record *model.FlowRecord) {
	if g == nil {
		return
	}
	if g.city != nil {
		if src := net.ParseIP(record.SrcIP); src != nil {
			if city, err := g.city.City(src); err == nil {
				if record.SrcCountry == "" {
					record.SrcCountry = city.Country.IsoCode
				}
			}
		}
		if dst := net.ParseIP(record.DstIP); dst != nil {
			if city, err := g.city.City(dst); err == nil {
				if record.DstCountry == "" {
					record.DstCountry = city.Country.IsoCode
				}
			}
		}
	}
	if g.asn != nil {
		if src := net.ParseIP(record.SrcIP); src != nil {
			if asn, err := g.asn.ASN(src); err == nil && record.SrcASN == 0 {
				record.SrcASN = uint32(asn.AutonomousSystemNumber)
			}
		}
		if dst := net.ParseIP(record.DstIP); dst != nil {
			if asn, err := g.asn.ASN(dst); err == nil && record.DstASN == 0 {
				record.DstASN = uint32(asn.AutonomousSystemNumber)
			}
		}
	}
}

type InventoryEnricher struct {
	index *inventory.Index
}

func NewInventoryEnricher(index *inventory.Index) *InventoryEnricher {
	return &InventoryEnricher{index: index}
}

func (i *InventoryEnricher) Enrich(_ context.Context, record *model.FlowRecord) {
	if i == nil || i.index == nil {
		return
	}
	if src := i.index.Match(record.SrcIP); src != nil {
		record.SrcHostname = src.Hostname
		record.SrcService = src.Service
		record.SrcEnvironment = src.Environment
		record.SrcOwnerTeam = src.OwnerTeam
	}
	if dst := i.index.Match(record.DstIP); dst != nil {
		record.DstHostname = dst.Hostname
		record.DstService = dst.Service
		record.DstEnvironment = dst.Environment
		record.DstOwnerTeam = dst.OwnerTeam
	}
}

type ReverseDNSEnricher struct {
	ttl   time.Duration
	mu    sync.RWMutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	value     string
	expiresAt time.Time
}

func NewReverseDNSEnricher(ttl time.Duration) *ReverseDNSEnricher {
	return &ReverseDNSEnricher{ttl: ttl, cache: map[string]cacheEntry{}}
}

func (r *ReverseDNSEnricher) Enrich(_ context.Context, record *model.FlowRecord) {
	if record.SrcHostname == "" {
		record.SrcHostname = r.lookup(record.SrcIP)
	}
	if record.DstHostname == "" {
		record.DstHostname = r.lookup(record.DstIP)
	}
}

func (r *ReverseDNSEnricher) lookup(ip string) string {
	now := time.Now()
	r.mu.RLock()
	if found, ok := r.cache[ip]; ok && found.expiresAt.After(now) {
		r.mu.RUnlock()
		return found.value
	}
	r.mu.RUnlock()
	name := ""
	names, err := net.LookupAddr(ip)
	if err == nil && len(names) > 0 {
		name = strings.TrimSuffix(names[0], ".")
	}
	r.mu.Lock()
	r.cache[ip] = cacheEntry{value: name, expiresAt: now.Add(r.ttl)}
	r.mu.Unlock()
	return name
}
