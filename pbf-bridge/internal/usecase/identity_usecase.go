package usecase

import (
	"bytes"
	"fmt"
	"text/template"

	"pbf-bridge/internal/domain"
)

type identityUseCase struct {
	printerRepo  domain.PrinterRepository
	templateRepo domain.TemplateRepository
	templatePath string
	pendingUC    domain.PendingUseCase
}

func NewIdentityUseCase(pr domain.PrinterRepository, tr domain.TemplateRepository, puc domain.PendingUseCase, path string) domain.IdentityUseCase {
	return &identityUseCase{
		printerRepo:  pr,
		templateRepo: tr,
		templatePath: path,
		pendingUC:    puc,
	}
}

func (uc *identityUseCase) ProcessIdentityLabels(payload domain.PrintIdentityPayload, isRetry bool) error {
	rawTemplate, err := uc.templateRepo.ReadTemplate(uc.templatePath)
	if err != nil {
		return fmt.Errorf("failed to read identity template file: %w", err)
	}

	// 2. Parsing text jadi Go Template
	tmpl, err := template.New("identity").Parse(rawTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse identity template: %w", err)
	}

	var finalZplBuffer bytes.Buffer

	// 3. Looping maju 2 langkah (chunking berpasangan)
	for i := 0; i < len(payload.Identities); i += 2 {
		label1 := payload.Identities[i]
		label1.QRCode = fmt.Sprintf("%s", label1.BatchNumber)

		label2 := domain.ProductIdentity{}
		if i+1 < len(payload.Identities) {
			label2 = payload.Identities[i+1]
			label2.QRCode = fmt.Sprintf("%s", label2.BatchNumber)
		}
		data := struct {
			Label1 domain.ProductIdentity
			Label2 domain.ProductIdentity
		}{
			Label1: label1,
			Label2: label2,
		}

		if err := tmpl.Execute(&finalZplBuffer, data); err != nil {
			return fmt.Errorf("failed to execute template at index %d: %w", i, err)
		}
	}

	err = uc.printerRepo.SendZPL(finalZplBuffer.String())
	if err != nil {
		if !isRetry {
			return uc.pendingUC.HandleFailedPrint("identity", payload, err)
		}
		return fmt.Errorf("printer connection failed during retry attempt: %w", err)
	}
	return nil
}
