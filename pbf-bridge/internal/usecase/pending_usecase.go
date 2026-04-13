package usecase

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"pbf-bridge/internal/domain"
)

type pendingUseCase struct {
	pendingRepo domain.PendingJobRepository
	shippingUC  domain.ShippingUseCase
	identityUC  domain.IdentityUseCase
}

func NewPendingUseCase(repo domain.PendingJobRepository) domain.PendingUseCase {
	return &pendingUseCase{
		pendingRepo: repo,
	}
}

func (uc *pendingUseCase) SetUsecases(shipping domain.ShippingUseCase, identity domain.IdentityUseCase) {
	uc.shippingUC = shipping
	uc.identityUC = identity
}

func (uc *pendingUseCase) HandleFailedPrint(jobType string, payload interface{}, originalErr error) error {
	saveErr := uc.pendingRepo.SaveFailedJob(jobType, payload)
	if saveErr != nil {
		slog.Error("Failed to save backup JSON", "job_type", jobType, "error", saveErr.Error())
		return fmt.Errorf("CRITICAL: printer offline (%v) and failed to save backup JSON: %v", originalErr, saveErr)
	}

	slog.Warn("Printer offline, data saved to pending folder", "job_type", jobType)
	return fmt.Errorf("printer offline, data successfully saved to pending folder: %w", originalErr)
}

func (uc *pendingUseCase) RetryAllPending() error {
	files, err := uc.pendingRepo.ListPendingFiles()
	if err != nil {
		slog.Error("Failed to list pending files", "error", err.Error())
		return fmt.Errorf("failed to list pending files: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	slog.Info("Starting retry for pending files", "total_files", len(files))

	for _, file := range files {
		jobType, rawData, err := uc.pendingRepo.ReadPendingFile(file)
		if err != nil {
			slog.Error("Failed to read pending file, skipping", "file", file, "error", err.Error())
			continue
		}

		var processErr error

		switch jobType {
		case "shipping":
			var payload domain.PrintShippingPayload
			if err := json.Unmarshal(rawData, &payload); err != nil {
				slog.Error("Failed to parse shipping JSON", "file", file, "error", err.Error())
				continue
			}
			// Add 'true' here
			processErr = uc.shippingUC.ProcessShippingLabels(payload, true)

		case "identity":
			var payload domain.PrintIdentityPayload
			if err := json.Unmarshal(rawData, &payload); err != nil {
				slog.Error("Failed to parse identity JSON", "file", file, "error", err.Error())
				continue
			}
			// Add 'true' here
			processErr = uc.identityUC.ProcessIdentityLabels(payload, true)

		default:
			slog.Warn("Unknown job type, skipping file", "file", file, "job_type", jobType)
			continue
		}

		if processErr == nil {
			slog.Info("Successfully retried job, deleting file", "file", file, "job_type", jobType)
			uc.pendingRepo.DeletePendingFile(file)
		} else {
			slog.Error("Retry failed", "file", file, "job_type", jobType, "error", processErr.Error())
		}
	}
	return nil
}
