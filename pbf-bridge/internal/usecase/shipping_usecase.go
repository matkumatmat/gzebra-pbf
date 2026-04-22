package usecase

import (
	"bytes"
	"fmt"
	"text/template"

	"pbf-bridge/internal/domain"
)

type shippingUseCase struct {
	printerRepo  domain.PrinterRepository
	templateRepo domain.TemplateRepository
	templatePath string
	pendingUC    domain.PendingUseCase
}

func NewShippingUseCase(pr domain.PrinterRepository, tr domain.TemplateRepository, puc domain.PendingUseCase, path string) domain.ShippingUseCase {
	return &shippingUseCase{
		printerRepo:  pr,
		templateRepo: tr,
		templatePath: path,
		pendingUC:    puc,
	}
}

func (uc *shippingUseCase) ProcessShippingLabels(payload domain.PrintShippingPayload, isRetry bool) error {
	rawTemplate, err := uc.templateRepo.ReadTemplate(uc.templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New("shipping").Parse(rawTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Iterate over each box and send them individually to avoid buffer overflow
	for _, box := range payload.Boxes {
		var labelBuffer bytes.Buffer

		data := struct {
			Recipient domain.ShippingRecipient
			TotalBox  int
			Box       domain.ShippingBox
		}{
			Recipient: payload.Recipient,
			TotalBox:  payload.TotalBox,
			Box:       box,
		}

		if err := tmpl.Execute(&labelBuffer, data); err != nil {
			return fmt.Errorf("failed to execute template for box %d: %w", box.CurrentBox, err)
		}

		// Send immediately inside the loop
		err = uc.printerRepo.SendZPL(labelBuffer.String())
		if err != nil {
			if !isRetry {
				return uc.pendingUC.HandleFailedPrint("shipping", payload, err)
			}
			return fmt.Errorf("printer connection failed during retry attempt: %w", err)
		}
	}

	return nil
}
