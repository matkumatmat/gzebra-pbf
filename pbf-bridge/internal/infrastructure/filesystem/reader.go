package filesystem

import (
	"os"
)

type fileTemplateRepo struct{}

func NewFileTemplateRepository() *fileTemplateRepo {
	return &fileTemplateRepo{}
}

func (r *fileTemplateRepo) ReadTemplate(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
