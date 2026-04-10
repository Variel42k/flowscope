package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/flowscope/flowscope/internal/config"
	"github.com/flowscope/flowscope/internal/decoder"
	"github.com/flowscope/flowscope/internal/enrich"
	"github.com/flowscope/flowscope/internal/inventory"
	"github.com/flowscope/flowscope/internal/model"
	"github.com/flowscope/flowscope/internal/normalize"
	"github.com/flowscope/flowscope/internal/storage"
)

type packet struct {
	sourceType string
	exporterIP string
	payload    []byte
}

func main() {
	cfg := config.Load()
	db, err := storage.Open(cfg.ClickHouseDSN)
	if err != nil {
		log.Fatalf("open storage: %v", err)
	}
	defer db.Close()
	repo := storage.NewRepository(db, cfg.RetentionDays)
	if err := repo.ApplyRetention(context.Background(), cfg.RetentionDays); err != nil {
		log.Printf("retention apply warning: %v", err)
	}

	var enrichers []enrich.Enricher
	var inv *inventory.Index
	if cfg.EnableInventory {
		inv, err = inventory.Load(cfg.InventoryPath, cfg.InventoryFormat)
		if err != nil {
			log.Printf("inventory load failed: %v", err)
		} else {
			enrichers = append(enrichers, enrich.NewInventoryEnricher(inv))
			if err := repo.ImportInventory(context.Background(), inv.Assets()); err != nil {
				log.Printf("inventory initial import failed: %v", err)
			}
		}
	}
	if cfg.EnableGeoIP {
		geo, err := enrich.NewGeoIPEnricher(cfg.GeoIPCityPath, cfg.GeoIPASNPath)
		if err != nil {
			log.Printf("geoip disabled: %v", err)
		} else {
			enrichers = append(enrichers, geo)
		}
	}
	if cfg.EnableRDNS {
		enrichers = append(enrichers, enrich.NewReverseDNSEnricher(cfg.RDNSCacheTTL))
	}
	if cfg.EnableInterfaceAlias {
		ifAlias, err := enrich.NewInterfaceAliasEnricher(cfg.InterfaceAliasPath)
		if err != nil {
			log.Printf("interface alias enrichment disabled: %v", err)
		} else {
			enrichers = append(enrichers, ifAlias)
		}
	}
	pipeline := enrich.NewComposite(enrichers...)
	dec := decoder.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	g, gctx := errgroup.WithContext(ctx)

	packets := make(chan packet, 2048)
	g.Go(func() error { return listenUDP(gctx, cfg.CollectorNetFlowV5, "netflow_v5", packets) })
	g.Go(func() error { return listenUDP(gctx, cfg.CollectorNetFlowV9, "netflow_v9", packets) })
	g.Go(func() error { return listenUDP(gctx, cfg.CollectorIPFIX, "ipfix", packets) })
	g.Go(func() error { return listenUDP(gctx, cfg.CollectorSFlow, "sflow", packets) })
	g.Go(func() error {
		defer close(packets)
		<-gctx.Done()
		return nil
	})
	g.Go(func() error {
		flushTicker := time.NewTicker(1 * time.Second)
		defer flushTicker.Stop()
		buffer := make([]model.FlowRecord, 0, 2048)
		flush := func() {
			if len(buffer) == 0 {
				return
			}
			if err := repo.InsertFlows(gctx, buffer); err != nil {
				log.Printf("insert batch failed: %v", err)
			}
			buffer = buffer[:0]
		}
		for {
			select {
			case <-gctx.Done():
				flush()
				return nil
			case pkt, ok := <-packets:
				if !ok {
					flush()
					return nil
				}
				records, err := dec.Decode(pkt.sourceType, pkt.payload, pkt.exporterIP)
				if err != nil {
					log.Printf("decode %s from %s failed: %v", pkt.sourceType, pkt.exporterIP, err)
					continue
				}
				now := time.Now().UTC()
				for _, rec := range records {
					normalize.Apply(&rec)
					pipeline.Enrich(gctx, &rec)
					buffer = append(buffer, rec)
					_ = repo.UpsertExporter(gctx, model.Exporter{
						ExporterID:     rec.ExporterID,
						ExporterIP:     rec.ExporterIP,
						ExporterName:   rec.ExporterName,
						ObservationID:  rec.ObservationID,
						FirstSeen:      rec.TimestampStart,
						LastSeen:       now,
						FlowsObserved:  1,
						LastSourceType: rec.SourceType,
					})
				}
				if len(buffer) >= 1000 {
					flush()
				}
			case <-flushTicker.C:
				flush()
			}
		}
	})

	log.Printf("collector listening: nf5=%s nf9=%s ipfix=%s sflow=%s", cfg.CollectorNetFlowV5, cfg.CollectorNetFlowV9, cfg.CollectorIPFIX, cfg.CollectorSFlow)
	if err := g.Wait(); err != nil {
		log.Fatalf("collector stopped with error: %v", err)
	}
}

func listenUDP(ctx context.Context, addr string, sourceType string, out chan<- packet) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	buffer := make([]byte, 65535)
	for {
		_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, remote, err := conn.ReadFrom(buffer)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				select {
				case <-ctx.Done():
					return nil
				default:
					continue
				}
			}
			return err
		}
		payload := make([]byte, n)
		copy(payload, buffer[:n])
		select {
		case <-ctx.Done():
			return nil
		case out <- packet{sourceType: sourceType, exporterIP: hostOnly(remote.String()), payload: payload}:
		}
	}
}

func hostOnly(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}
