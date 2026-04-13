# This Document containing guide for connecting Zebra Printer to devices

**Linux (Arch Linux)**
1. install libraries and drivers
    ```bash
    sudo pacman -S gnu-netcat cups
    ```
2. connecting devices with printer with same networks (LAN or USB), and try to check connectivity
    ```bash
    ip -a
    ```

3. set ip address manually
   ```bash
    sudo ip addr add 192.168.19.10/24 dev enp0s20f0u1c2
   ```
   OR
   ```bash
    sudo ip link set dev enp0s20f0u1c2 up
   ```
4. verify connection
   ```bash
    ping -c 4 192.168.19.5
   ```

5. Test Print
   ```bash
    echo "^XA^FO50,50^A0N,50,50^FDTest from Arch Linux!^FS^XZ" | nc -w 1 192.168.19.5 9100
   ```