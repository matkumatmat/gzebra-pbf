package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort           string
	PrinterIP            string
	PrinterPort          string
	PrinterTimeoutSec    int
	ShippingTemplatePath string
	IdentityTemplatePath string
	PendingJobPath       string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Can't load .env file, relying on environment variables", "error", err.Error())
	}

	timeout, _ := strconv.Atoi(os.Getenv("PRINTER_TIMEOUT_SEC"))
	if timeout == 0 {
		timeout = 5
	}

	pendingPath := os.Getenv("PENDING_JOB_PATH")
	if pendingPath == "" {
		pendingPath = "./data/pending"
	}

	return &Config{
		ServerPort:           os.Getenv("SERVER_PORT"),
		PrinterIP:            os.Getenv("PRINTER_IP"),
		PrinterPort:          os.Getenv("PRINTER_PORT"),
		PrinterTimeoutSec:    timeout,
		ShippingTemplatePath: os.Getenv("SHIPPING_TEMPLATE_PATH"),
		IdentityTemplatePath: os.Getenv("IDENTITY_TEMPLATE_PATH"),
		PendingJobPath:       pendingPath,
	}
}
