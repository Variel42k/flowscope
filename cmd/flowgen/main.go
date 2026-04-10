package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/flowscope/flowscope/internal/config"
	"github.com/flowscope/flowscope/internal/model"
	"github.com/flowscope/flowscope/internal/normalize"
	"github.com/flowscope/flowscope/internal/storage"
)

func main() {
	cfg := config.Load()
	db, err := storage.Open(cfg.ClickHouseDSN)
	if err != nil {
		log.Fatalf("flowgen: storage open failed: %v", err)
	}
	defer db.Close()
	repo := storage.NewRepository(db, cfg.RetentionDays)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	ctx := context.Background()
	for t := range ticker.C {
		batch := make([]model.FlowRecord, 0, 64)
		for i := 0; i < 64; i++ {
			proto := uint8([]int{6, 6, 6, 17, 1}[rng.Intn(5)])
			rec := model.FlowRecord{
				TimestampStart:  t.UTC(),
				TimestampEnd:    t.UTC().Add(time.Duration(rng.Intn(5)) * time.Second),
				ExporterID:      []string{"lab-router-a", "lab-router-b", "lab-sflow-a"}[rng.Intn(3)],
				ExporterIP:      []string{"172.16.0.1", "172.16.0.2", "172.16.0.3"}[rng.Intn(3)],
				ObservationID:   uint32(1 + rng.Intn(3)),
				SrcIP:           randomIP(rng, []string{"10.10.1", "10.10.2", "10.20.1", "192.0.2", "198.51.100"}),
				DstIP:           randomIP(rng, []string{"10.10.1", "10.10.2", "10.20.1", "192.0.2", "198.51.100"}),
				SrcPort:         uint16(10000 + rng.Intn(50000)),
				DstPort:         uint16([]int{443, 80, 53, 22, 5432, 6379, 3389, 8443}[rng.Intn(8)]),
				IPProtocol:      proto,
				Bytes:           uint64(400 + rng.Intn(500000)),
				Packets:         uint64(1 + rng.Intn(700)),
				InputInterface:  uint32([]int{10, 11, 20, 21, 100}[rng.Intn(5)]),
				OutputInterface: uint32([]int{10, 11, 20, 21, 100}[rng.Intn(5)]),
				SrcASN:          uint32([]int{64512, 64513, 13335, 15169}[rng.Intn(4)]),
				DstASN:          uint32([]int{64512, 64513, 13335, 15169}[rng.Intn(4)]),
				SrcCountry:      []string{"US", "DE", "FR", "RU", "SG"}[rng.Intn(5)],
				DstCountry:      []string{"US", "DE", "FR", "RU", "SG"}[rng.Intn(5)],
				TCPFlags:        uint16([]int{2, 16, 24, 17}[rng.Intn(4)]),
				SourceType:      []string{"netflow_v9", "ipfix", "sflow"}[rng.Intn(3)],
			}
			if i == 0 && t.Unix()%30 == 0 {
				rec.SrcIP = "10.10.1.99"
				rec.DstIP = randomIP(rng, []string{"198.51.100", "203.0.113"})
				rec.Bytes = 8_000_000
				rec.DstPort = 4444
			}
			normalize.Apply(&rec)
			batch = append(batch, rec)
		}
		if err := repo.InsertFlows(ctx, batch); err != nil {
			log.Printf("flowgen insert error: %v", err)
			continue
		}
		_ = repo.RunRollups(ctx, t.Add(-2*time.Minute).UTC().Truncate(time.Minute), t.Add(time.Minute).UTC().Truncate(time.Minute))
	}
}

func randomIP(rng *rand.Rand, prefixes []string) string {
	p := prefixes[rng.Intn(len(prefixes))]
	return p + "." + strconv.Itoa(1+rng.Intn(253))
}
