package printer

import (
	"fmt"
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
	// Menggunakan net.JoinHostPort menggantikan fmt.Sprintf
	address := net.JoinHostPort(p.ip, p.port)
	timeout := time.Duration(p.timeoutSec) * time.Second
	fmt.Println("🕵️ DEBUG: Go lagi nyoba nembak ke IP ->", address)

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return fmt.Errorf("koneksi printer gagal di %s: %w", address, err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(zpl))
	if err != nil {
		return fmt.Errorf("gagal mengirim ZPL byte: %w", err)
	}

	return nil
}
