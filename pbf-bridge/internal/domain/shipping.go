package domain

type ShippingProduct struct {
	Name  string `json:"name"`
	Qty   string `json:"qty"`
	Batch string `json:"batch"`
	Unit  string `json:"unit"`
}

type ShippingBox struct {
	CurrentBox  int               `json:"current_box"`
	Petugas     string            `json:"petugas"`
	Temperature string            `json:"temperature"`
	Products    []ShippingProduct `json:"products"`
}

type ShippingRecipient struct {
	Customer     string `json:"customer"`
	Branch       string `json:"branch"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	Contact      string `json:"contact"`
}

// PrintShippingPayload payload utama untuk pengiriman (WMS)
type PrintShippingPayload struct {
	Recipient ShippingRecipient `json:"recipient"`
	TotalBox  int               `json:"total_box"`
	Boxes     []ShippingBox     `json:"boxes"`
}

// ShippingUseCase kontrak untuk business logic label pengiriman
type ShippingUseCase interface {
	ProcessShippingLabels(payload PrintShippingPayload, isRetry bool) error
}
