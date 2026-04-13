# pbf-bridge

A lightweight bridge service that connects web applications (Google AppSheet / Apps Script) to Zebra label printers via ZPL over TCP/IP.

---

## Getting Started

### 1. Clone & Enter Directory
```bash
git clone https://github.com/matkumatmat/gzebra-pbf.git
cd pbf-bridge
```

### 2. Configure Environment
Copy `.env.example` to `.env` and adjust the values:
```bash
cp .env.example .env
```

```env
SERVER_PORT=8080
PRINTER_IP=192.168.19.5       # Zebra printer IP
PRINTER_PORT=9112              # Printer port (Zebra default: 9100 or 9112)
PRINTER_TIMEOUT_SEC=5

PENDING_JOB_PATH=./data/pending
SHIPPING_TEMPLATE_PATH=./templates/shipping.txt
IDENTITY_TEMPLATE_PATH=./templates/product.txt
```

### 3. Run the Binary

**Linux:**
```bash
./bin/pbf-bridge-linux
```

**Windows** — open CMD or PowerShell from inside the repo folder:
```
bin\pbf-bridge-windows.exe
```

> ⚠️ **Important:** Always run from inside the repo folder. Do not double-click the `.exe` directly from Explorer, as it may not resolve the `.env` and `templates/` paths correctly.

---

## Endpoints

Base URL: `http://localhost:8080`

### POST `/print/shipping`
Print a shipping label per box.

**Payload:**
```json
{
  "recipient": {
    "customer": "Klinik Sehat Selalu",
    "branch": "Cabang Pasteur",
    "address_line_1": "Jl. Dr. Djunjunan No. 123",
    "address_line_2": "Kota Bandung, 40162",
    "contact": "Bpk. Budi - 08123456789"
  },
  "total_box": 2,
  "boxes": [
    {
      "current_box": 1,
      "petugas": "Kaye",
      "temperature": "2-8 °C",
      "products": [
        { "name": "Vaksin A", "qty": "10" }
      ]
    }
  ]
}
```

---

### POST `/print/identity`
Print product identity labels. The template fits 2 labels per page — items are automatically paired.

**Payload:**
```json
{
  "identities": [
    {
      "product_code": "VAK-001",
      "product_name_1": "VAKSIN INFLUENZA",
      "product_name_2": "STRAIN SELATAN 2026",
      "batch_number": "BCH-998877",
      "allocation": "KOTA BANDUNG",
      "mfg_date": "10/03/2026",
      "exp_date": "10/03/2028",
      "receive_date": "12/03/2026",
      "qr_code": "VAK-001|BCH-998877|10/03/2028"
    }
  ]
}
```

---

## Project Structure

```
pbf-bridge/
├── bin/
│   ├── pbf-bridge-linux        # Linux executable
│   └── pbf-bridge-windows.exe  # Windows executable
├── templates/
│   ├── shipping.txt            # ZPL shipping label template (editable)
│   └── product.txt             # ZPL identity label template (editable)
├── data/
│   └── pending/                # Auto-generated — stores failed jobs for retry
├── .env                        # Configuration (create from .env.example)
├── .env.example
└── build.sh                    # Cross-platform build script
```

---

## Editing Templates

ZPL templates are located in the `templates/` folder. Edit the `.txt` files directly — no rebuild required. Restart the service after making changes.

---

## Pending & Retry Mechanism

If the printer is offline when a request comes in, the payload is automatically saved to `data/pending/` as a JSON file. The service retries all pending jobs every **1 minute** once the printer is back online.

---

## Build from Source

Requires Go to be installed.

```bash
chmod +x build.sh
./build.sh
```

Output binaries will be placed in the `bin/` folder.

---

## Local Testing (without a printer)

Simulate a printer using `netcat`:

**Linux / Mac:**
```bash
nc -l 54321
```

**Windows** (using ncat from Nmap or WSL):
```bash
ncat -l 54321
```

Then update `.env`:
```env
PRINTER_IP=127.0.0.1
PRINTER_PORT=54321
```

---

## Changelog

| Version | Notes |
|---------|-------|
| v0.1.0 | Initial release — shipping & identity label printing, pending retry mechanism |