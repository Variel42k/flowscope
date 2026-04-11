package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/flowscope/flowscope/internal/config"
	"github.com/flowscope/flowscope/internal/rollup"
	"github.com/flowscope/flowscope/internal/storage"
)

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
	if err := repo.EnsureDefaultAlertRules(context.Background(), cfg.AdminUser); err != nil {
		log.Printf("alert rule bootstrap warning: %v", err)
	}
	worker := rollup.NewWorker(repo, cfg.RollupInterval)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log.Printf("worker running, interval=%s", cfg.RollupInterval)
	worker.Run(ctx)
}
