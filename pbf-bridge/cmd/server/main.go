package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"pbf-bridge/internal/config"
	delivery "pbf-bridge/internal/delivery/http"
	"pbf-bridge/internal/infrastructure/filesystem"
	"pbf-bridge/internal/infrastructure/printer"
	"pbf-bridge/internal/usecase"
	"pbf-bridge/pkg/cron"
	"pbf-bridge/pkg/logger"
)

func main() {
	// 1. Init Logger
	logger.Init()

	// 2. Load Config
	cfg := config.LoadConfig()

	// 3. Setup Context for Graceful Shutdown of background workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 4. Init Repositories (Infrastructure Layer)
	printerRepo := printer.NewBidirectionalPrinter(cfg.PrinterIP, cfg.PrinterPort, cfg.PrinterTimeoutSec)
	templateRepo := filesystem.NewFileTemplateRepository()
	pendingRepo := filesystem.NewPendingJobRepository(cfg.PendingJobPath)

	// 5. Init Pending Usecase (Without main usecases to prevent circular dependency)
	pendingUC := usecase.NewPendingUseCase(pendingRepo)

	// 6. Init Main Usecases (Injecting pendingUC)
	shippingUC := usecase.NewShippingUseCase(printerRepo, templateRepo, pendingUC, cfg.ShippingTemplatePath)
	identityUC := usecase.NewIdentityUseCase(printerRepo, templateRepo, pendingUC, cfg.IdentityTemplatePath)

	// 7. Resolve Circular Dependency
	pendingUC.SetUsecases(shippingUC, identityUC)

	// 8. Init and Start Internal Cron Worker
	// Setting interval to 30 minutes. You can extract this to .env/config later if needed.
	retryWorker := cron.NewWorker(pendingUC, 30)
	go retryWorker.Start(ctx)

	// 9. Init Handler (Delivery Layer)
	printHandler := delivery.NewPrintHandler(shippingUC, identityUC)

	// 10. Setup HTTP Router
	mux := http.NewServeMux()
	mux.HandleFunc("/print/shipping", printHandler.PrintShipping)
	mux.HandleFunc("/print/identity", printHandler.PrintIdentity)

	// 11. Run Server
	address := fmt.Sprintf(":%s", cfg.ServerPort)
	slog.Info("Bridge system running", "port", cfg.ServerPort)

	if err := http.ListenAndServe(address, mux); err != nil {
		slog.Error("Server stopped", "error", err.Error())
		os.Exit(1)
	}
}
