import socket
import os

# Bikin folder debug kalo belum ada
os.makedirs("debug", exist_ok=True)

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
# Pake reuse address biar ga error 'address already in use' pas di-restart
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.bind(('127.0.0.1', 9999))
s.listen(1)

print("🖨️  Fake Zebra Printer nyala di 127.0.0.1:9999...")

while True:
    conn, addr = s.accept()
    data = conn.recv(1024)
    
    if b"~HS" in data:
        print("-> Go nanya status (~HS). Ngebales: KERTAS AMAN!")
        # Balasan Zebra Sehat (Index Kertas=0, Pause=0, Head=0, Ribbon=0)
        conn.sendall(b"\x02030,0,0,087,0,0,0,0,0,0,0,0\x03\x02000,0,0,0,0,0,0,0,0,0,0,0\x03")
    
    # Nungguin Go ngirim ZPL aslinya
    zpl = conn.recv(4096)
    if zpl:
        with open("debug/debug.txt", "ab") as f:
            f.write(zpl)
            f.write(b"\n--- END OF LABEL ---\n")
        print("-> ZPL Diterima & di-save ke debug/debug.txt!")
    
    conn.close()