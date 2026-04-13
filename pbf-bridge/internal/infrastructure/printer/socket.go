package printer

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

type socketPrinter struct {
	ip         string
	port       string
	timeoutSec int
}

func NewSocketPrinter(ip, port string, timeoutSec int) *socketPrinter {
	return &socketPrinter{
		ip:         ip,
		port:       port,
		timeoutSec: timeoutSec,
	}
}

func (p *socketPrinter) SendZPL(zpl string) error {
	address := net.JoinHostPort(p.ip, p.port)
	timeout := time.Duration(p.timeoutSec) * time.Second
	// fmt.Println("DEBUG: request to :", address)
	slog.Info("Sending ZPL request", slog.String("address", address))

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return fmt.Errorf("printer connection failed on: %s: %w", address, err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(zpl))
	if err != nil {
		return fmt.Errorf("failed to send ZPL bytes on: %s: %w", address, err)
	}

	return nil
}
