package printer

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type bidirectionalPrinter struct {
	ip         string
	port       string
	timeoutSec int
}

func NewBidirectionalPrinter(ip, port string, timeoutSec int) *bidirectionalPrinter {
	return &bidirectionalPrinter{
		ip:         ip,
		port:       port,
		timeoutSec: timeoutSec,
	}
}

func (p *bidirectionalPrinter) SendZPL(zpl string) error {
	address := net.JoinHostPort(p.ip, p.port)
	timeout := time.Duration(p.timeoutSec) * time.Second

	// 1. SYNC
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return fmt.Errorf("koneksi ke printer mati/putus: %w", err)
	}
	defer conn.Close()

	// 2. ACK
	_, err = conn.Write([]byte("~HS"))
	if err != nil {
		return fmt.Errorf("gagal kirim request status ke printer: %w", err)
	}

	// 3. LISTEN PRINTER RESPONSES (Status Hardware)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("printer connected but got fail to read response: %w", err)
	}

	// 4. PARSING RESPONSE PRINTER
	rawResponse := string(buffer[:n])
	if err := p.checkHardwareStatus(rawResponse); err != nil {
		return fmt.Errorf("HARDWARE ERROR: %w", err)
	}

	conn.SetDeadline(time.Time{})
	_, err = conn.Write([]byte(zpl))
	if err != nil {
		return fmt.Errorf("status aman, tapi gagal mengirim ZPL: %w", err)
	}

	return nil
}

// HANDLER PARSING PRINTER ZEBRA ORDER STATUS HARDWARE
func (p *bidirectionalPrinter) checkHardwareStatus(rawResponse string) error {
	cleanResp := strings.ReplaceAll(rawResponse, "\x02", "")
	blocks := strings.Split(cleanResp, "\x03")

	if len(blocks) < 2 {
		return fmt.Errorf("format balasan tidak valid, bukan standar Zebra")
	}
	// BLOK 1: Cek Kertas & Status Pause
	// Format string Zebra: aaa,b,c,dddd,e,f,g,h,i,j,k,l
	// Index 1 = Paper Out (1 kalau habis)
	// Index 2 = Pause (1 kalau lagi di-pause)
	block1 := strings.Split(blocks[0], ",")
	if len(block1) >= 3 {
		if strings.TrimSpace(block1[1]) == "1" {
			return fmt.Errorf("KERTAS LABEL HABIS")
		}
		if strings.TrimSpace(block1[2]) == "1" {
			return fmt.Errorf("PRINTER SEDANG DI-PAUSE")
		}
	}

	// BLOK 2: Cek Tutup Head & Ribbon (Tinta)
	// Format string Zebra: aaa,b,c,d,e,f,g,h,i,j,k,l
	// Index 4 = Head Open (1 kalau tutup terbuka)
	// Index 5 = Ribbon Out (1 kalau tinta habis)
	block2 := strings.Split(blocks[1], ",")
	if len(block2) >= 6 {
		if strings.TrimSpace(block2[4]) == "1" {
			return fmt.Errorf("TUTUP PRINTER (HEAD) TERBUKA")
		}
		if strings.TrimSpace(block2[5]) == "1" {
			return fmt.Errorf("RIBBON TINTA HABIS")
		}
	}
	return nil
}
