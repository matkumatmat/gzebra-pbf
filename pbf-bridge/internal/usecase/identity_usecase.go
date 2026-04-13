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
		return fmt.Errorf("gagal membaca file template identitas: %w", err)
	}

	// 2. Parsing text jadi Go Template
	tmpl, err := template.New("identity").Parse(rawTemplate)
	if err != nil {
		return fmt.Errorf("gagal parsing template identitas: %w", err)
	}

	var finalZplBuffer bytes.Buffer

	// 3. Looping maju 2 langkah (chunking berpasangan)
	for i := 0; i < len(payload.Identities); i += 2 {
		label1 := payload.Identities[i]

		// MAGIC QR CODE 1: Rakit otomatis string QR Code (Kode|Batch|Exp)
		// label1.QRCode = fmt.Sprintf("%s|%s|%s", label1.ProductCode, label1.BatchNumber, label1.ExpDate)
		// test use batch
		label1.QRCode = fmt.Sprintf("%s", label1.BatchNumber)

		// Bikin Label 2 default kosong dulu
		label2 := domain.ProductIdentity{}

		// Kalau index berikutnya masih ada di array, masukin ke Label 2
		if i+1 < len(payload.Identities) {
			label2 = payload.Identities[i+1]

			// MAGIC QR CODE 2: Rakit otomatis buat label bawahnya
			// label2.QRCode = fmt.Sprintf("%s|%s|%s", label2.ProductCode, label2.BatchNumber, label2.ExpDate)
			// test use batch
			label2.QRCode = fmt.Sprintf("%s", label2.BatchNumber)
		}

		// Bungkus berdua buat dilempar ke template
		data := struct {
			Label1 domain.ProductIdentity
			Label2 domain.ProductIdentity
		}{
			Label1: label1,
			Label2: label2,
		}

		if err := tmpl.Execute(&finalZplBuffer, data); err != nil {
			return fmt.Errorf("gagal render template identitas index %d: %w", i, err)
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
