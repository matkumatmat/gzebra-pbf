package domain

type ProductIdentity struct {
	ProductCode  string `json:"product_code"`
	ProductName1 string `json:"product_name_1"`
	ProductName2 string `json:"product_name_2"`
	BatchNumber  string `json:"batch_number"`
	Allocation   string `json:"allocation"`
	MfgDate      string `json:"mfg_date"`
	ExpDate      string `json:"exp_date"`
	ReceiveDate  string `json:"receive_date"`
	QRCode       string `json:"qr_code"`
}

type PrintIdentityPayload struct {
	Identities []ProductIdentity `json:"identities"`
}

type IdentityUseCase interface {
	ProcessIdentityLabels(payload PrintIdentityPayload, isRetry bool) error
}
