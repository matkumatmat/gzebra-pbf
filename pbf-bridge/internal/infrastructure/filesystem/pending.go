package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type pendingJobRepo struct {
	basePath string
}

func NewPendingJobRepository(basePath string) *pendingJobRepo {
	return &pendingJobRepo{
		basePath: basePath,
	}
}

func (r *pendingJobRepo) SaveFailedJob(jobType string, payload interface{}) error {
	// 1. Pastikan foldernya ada (kalau belum ada, dibikinin otomatis)
	if err := os.MkdirAll(r.basePath, os.ModePerm); err != nil {
		return fmt.Errorf("gagal membuat folder pending: %w", err)
	}

	// 2. Bikin nama file berdasarkan waktu & tipe job
	// Format: 2026-04-13_22-57-00_shipping_12345.json
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	nano := time.Now().UnixNano() % 100000 // Biar unik kalau gagalnya barengan
	filename := fmt.Sprintf("%s_%s_%d.json", timestamp, jobType, nano)
	fullPath := filepath.Join(r.basePath, filename)

	// 3. Marshal payload struct kembali menjadi JSON text (Pretty print)
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("gagal marshal JSON backup: %w", err)
	}

	// 4. Tulis file-nya
	if err := os.WriteFile(fullPath, jsonData, 0644); err != nil {
		return fmt.Errorf("gagal menulis file pending: %w", err)
	}

	return nil
}

func (r *pendingJobRepo) ListPendingFiles() ([]string, error) {
	entries, err := os.ReadDir(r.basePath)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			files = append(files, e.Name())
		}
	}
	return files, nil
}

func (r *pendingJobRepo) DeletePendingFile(filename string) error {
	return os.Remove(filepath.Join(r.basePath, filename))
}

func (r *pendingJobRepo) ReadPendingFile(filename string) (string, []byte, error) {
	fullPath := filepath.Join(r.basePath, filename)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", nil, err
	}

	parts := strings.Split(filename, "_")
	jobType := "unknown"
	if len(parts) >= 3 {
		jobType = parts[2]
	}

	return jobType, data, nil
}
