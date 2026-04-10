package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/flowscope/flowscope/internal/config"
	"github.com/flowscope/flowscope/internal/inventory"
	"github.com/flowscope/flowscope/internal/model"
	"github.com/flowscope/flowscope/internal/normalize"
	"github.com/flowscope/flowscope/internal/storage"
)

func main() {
	cfg := config.Load()
	db, err := storage.Open(cfg.ClickHouseDSN)
	if err != nil {
		log.Fatalf("seed: storage open failed: %v", err)
	}
	defer db.Close()
	repo := storage.NewRepository(db, cfg.RetentionDays)
	ctx := context.Background()

	if cfg.EnableInventory {
		idx, err := inventory.Load(cfg.InventoryPath, cfg.InventoryFormat)
		if err == nil {
			if err := repo.ImportInventory(ctx, idx.Assets()); err != nil {
				log.Printf("seed: inventory import failed: %v", err)
			}
		}
	}

	rng := rand.New(rand.NewSource(42))
	now := time.Now().UTC().Truncate(time.Minute)
	total := 0
	batch := make([]model.FlowRecord, 0, 2000)
	exporters := []model.Exporter{
		{ExporterID: "lab-router-a", ExporterIP: "172.16.0.1", ExporterName: "lab-router-a", ObservationID: 1, LastSourceType: "netflow_v9"},
		{ExporterID: "lab-router-b", ExporterIP: "172.16.0.2", ExporterName: "lab-router-b", ObservationID: 2, LastSourceType: "ipfix"},
		{ExporterID: "lab-sflow-a", ExporterIP: "172.16.0.3", ExporterName: "lab-sflow-a", ObservationID: 3, LastSourceType: "sflow"},
	}
	for i := 0; i < 16000; i++ {
		exp := exporters[i%len(exporters)]
		t := now.Add(-time.Duration(rng.Intn(360)) * time.Minute)
		src := randomIP(rng, []string{"10.10.1", "10.10.2", "10.20.1", "192.0.2", "198.51.100"})
		dst := randomIP(rng, []string{"10.10.1", "10.10.2", "10.20.1", "192.0.2", "198.51.100"})
		for src == dst {
			dst = randomIP(rng, []string{"10.10.1", "10.10.2", "10.20.1", "192.0.2", "198.51.100"})
		}
		proto := uint8([]int{6, 6, 6, 17, 1}[rng.Intn(5)])
		port := uint16([]int{443, 80, 53, 22, 5432, 6379, 3389}[rng.Intn(7)])
		bytes := uint64(400 + rng.Intn(900000))
		if i%113 == 0 {
			bytes *= 8
		}
		rec := model.FlowRecord{
			TimestampStart:  t,
			TimestampEnd:    t.Add(time.Duration(rng.Intn(20)) * time.Second),
			ExporterID:      exp.ExporterID,
			ExporterIP:      exp.ExporterIP,
			ExporterName:    exp.ExporterName,
			ObservationID:   exp.ObservationID,
			SrcIP:           src,
			DstIP:           dst,
			SrcPort:         uint16(10000 + rng.Intn(50000)),
			DstPort:         port,
			IPProtocol:      proto,
			Bytes:           bytes,
			Packets:         uint64(1 + rng.Intn(1600)),
			InputInterface:  uint32([]int{10, 11, 20, 21, 100}[rng.Intn(5)]),
			OutputInterface: uint32([]int{10, 11, 20, 21, 100}[rng.Intn(5)]),
			SrcASN:          uint32([]int{64512, 64513, 64514, 13335, 15169}[rng.Intn(5)]),
			DstASN:          uint32([]int{64512, 64513, 64514, 13335, 15169}[rng.Intn(5)]),
			SrcCountry:      []string{"US", "DE", "FR", "RU", "SG"}[rng.Intn(5)],
			DstCountry:      []string{"US", "DE", "FR", "RU", "SG"}[rng.Intn(5)],
			TCPFlags:        uint16([]int{2, 16, 24, 17}[rng.Intn(4)]),
			SourceType:      exp.LastSourceType,
		}
		normalize.Apply(&rec)
		batch = append(batch, rec)
		total++
		if len(batch) >= 2000 {
			if err := repo.InsertFlows(ctx, batch); err != nil {
				log.Fatalf("seed insert failed: %v", err)
			}
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		if err := repo.InsertFlows(ctx, batch); err != nil {
			log.Fatalf("seed insert failed: %v", err)
		}
	}
	for _, exp := range exporters {
		exp.FirstSeen = now.Add(-6 * time.Hour)
		exp.LastSeen = now
		exp.FlowsObserved = uint64(total / len(exporters))
		_ = repo.UpsertExporter(ctx, exp)
	}
	if err := repo.RunRollups(ctx, now.Add(-6*time.Hour), now.Add(1*time.Minute)); err != nil {
		log.Printf("seed rollup warning: %v", err)
	}
	log.Printf("seed complete: inserted %d flow records", total)
	if os.Getenv("FLOWSCOPE_SEED_EXIT_SUCCESS") == "1" {
		os.Exit(0)
	}
}

func randomIP(rng *rand.Rand, prefixes []string) string {
	p := prefixes[rng.Intn(len(prefixes))]
	return p + "." + itoa(1+rng.Intn(253))
}

func itoa(v int) string {
	return strconv.Itoa(v)
}
