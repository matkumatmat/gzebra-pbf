import socket
import os

os.makedirs("debug", exist_ok=True)

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(('127.0.0.1', 9999))
s.listen(1)

print("Fake Zebra Printer nyala di 127.0.0.1:9999...")

while True:
    conn, addr = s.accept()
    print(f"-> Koneksi dari {addr}")

    raw = b""
    conn.settimeout(2.0)
    try:
        while True:
            chunk = conn.recv(4096)
            if not chunk:
                break
            raw += chunk
            
            # TAMBAHAN: Kalau udah nerima ~HS, langsung break biar gak nunggu timeout 2 detik!
            if b"~HS" in raw:
                break
                
    except socket.timeout:
        pass

    print(f"-> Total raw ({len(raw)} bytes), preview: {raw[:30]}")

    if b"~HS" in raw:
        print("-> Ada ~HS, ngebales status AMAN")
        conn.sendall(b"\x02030,0,0,087,0,0,0,0,0,0,0,0\x03\x02000,0,0,0,0,0,0,0,0,0,0,0\x03")

        # Sekarang tunggu ZPL aslinya
        zpl = b""
        conn.settimeout(3.0)
        try:
            while True:
                chunk = conn.recv(4096)
                if not chunk:
                    break
                zpl += chunk
        except socket.timeout:
            pass

        print(f"-> ZPL ({len(zpl)} bytes)")
        if zpl:
            with open("/home/k/go-lang/label-server/debug/debug.txt", "wb") as f:
                f.write(zpl)
            print("-> Saved!")
    else:
        print(f"-> Ga ada ~HS, skip. Data: {raw[:50]}")

    conn.close()