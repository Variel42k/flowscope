package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flowscope/flowscope/internal/api"
	"github.com/flowscope/flowscope/internal/auth"
	"github.com/flowscope/flowscope/internal/config"
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
	authMgr := auth.NewManager(cfg.JWTSecret)
	server := api.NewServer(cfg, repo, authMgr)
	httpServer := &http.Server{Addr: cfg.HTTPAddr, Handler: server.Handler(), ReadHeaderTimeout: 5 * time.Second}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = httpServer.Shutdown(shutdownCtx)
	}()
	log.Printf("api listening on %s", cfg.HTTPAddr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("api failed: %v", err)
	}
}
