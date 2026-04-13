package domain

type PrinterRepository interface {
	SendZPL(zpl string) error
}

type TemplateRepository interface {
	ReadTemplate(path string) (string, error)
}
