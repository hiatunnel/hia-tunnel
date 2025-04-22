# Hia‑Tunnel

A **light‑weight, QUIC‑based encrypted tunnel** that hides inside HTTP/3 TLS fingerprints, adds a second Noise _XK_ layer (with a pre‑shared key), multiplexes any number of TCP streams through one UDP flow, and gives you instant **SOCKS proxy + static port‑forwarding** — all driven by interactive shell menus.

---

## Features

- **QUIC-based encrypted tunnel** for privacy and performance.
- **HTTP/3 TLS camouflage** to blend in with regular web traffic.
- **Noise_XK (PSK) double encryption** for enhanced security.
- **Multiplex multiple TCP streams** over a single UDP connection.
- **SOCKS proxy & static port forwarding**, managed via interactive shell menus.
- **Easy one-command install** scripts for both server and client.

---

## One-command Install (Debian / Ubuntu)

### Server

```bash
curl -fsSL https://raw.githubusercontent.com/hiatunnel/hia-tunnel/main/scripts/install-server.sh | sudo bash
```
- The script **prompts for a listen port and PSK** (does **not** default to 443).
- Builds `hia-tunnel-server`, enables a systemd unit, and drops `/usr/local/bin/hia-menu` for daily management.

**Manage the server at any time:**
```bash
sudo hia-menu   # Add/delete port‑forward, change PSK or port, restart
```

---

### Client

```bash
curl -fsSL https://raw.githubusercontent.com/hiatunnel/hia-tunnel/main/scripts/install-client.sh | bash
```
- Use the interactive menu to **add/delete/list server peers**, each with its own PSK.
- Every peer gets its own local SOCKS port (`127.0.0.1:1080`, `1081`, …).

**Start as a systemd service:**
```bash
sudo systemctl enable --now hia-tunnel-client
```

---

## Quick Start

1. **Install the server** on your VPS.
2. **Set up the client** on your local machine.
3. **Use the menus** to add peers and configure port forwards.
4. Enjoy secure, multiplexed, and stealthy connections!

---