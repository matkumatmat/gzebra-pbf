package cron

import (
	"context"
	"log/slog"
	"time"

	"pbf-bridge/internal/domain"
)

type Worker struct {
	pendingUC domain.PendingUseCase
	interval  time.Duration
}

func NewWorker(puc domain.PendingUseCase, intervalMinutes int) *Worker {
	return &Worker{
		pendingUC: puc,
		interval:  time.Duration(intervalMinutes) * time.Minute,
	}
}

func (w *Worker) Start(ctx context.Context) {
	slog.Info("Starting internal cron worker", "interval_minutes", w.interval.Minutes())

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run once immediately on start
	if err := w.pendingUC.RetryAllPending(); err != nil {
		slog.Error("Initial retry failed", "error", err.Error())
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("Internal cron worker stopping")
			return
		case <-ticker.C:
			slog.Info("Cron worker triggering retry process")
			if err := w.pendingUC.RetryAllPending(); err != nil {
				slog.Error("Periodic retry failed", "error", err.Error())
			}
		}
	}
}
