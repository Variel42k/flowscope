package rollup

import (
	"context"
	"log"
	"time"

	"github.com/flowscope/flowscope/internal/storage"
)

type Worker struct {
	repo     *storage.Repository
	interval time.Duration
}

func NewWorker(repo *storage.Repository, interval time.Duration) *Worker {
	return &Worker{repo: repo, interval: interval}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		if err := w.runOnce(ctx); err != nil {
			log.Printf("rollup run failed: %v", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *Worker) runOnce(ctx context.Context) error {
	to := time.Now().UTC().Truncate(time.Minute)
	from := to.Add(-10 * time.Minute)
	return w.repo.RunRollups(ctx, from, to)
}
