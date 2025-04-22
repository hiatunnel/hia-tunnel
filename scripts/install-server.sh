#!/usr/bin/env bash
set -euo pipefail

# ── 固定缺省端口 / PSK（可用环境变量覆盖） ────────────────
PORT="${PORT:-8443}"
PSK="${PSK:-$(openssl rand -hex 32)}"
echo -e "\n[Hia‑Tunnel] 监听端口 : $PORT"
echo -e "[Hia‑Tunnel] 预共享密钥 : $PSK\n"

# ── 安装官方 Go 1.22（避免 apt 默认 1.19） ────────────────
GO_URL="https://go.dev/dl/go1.19.8.linux-amd64.tar.gz"
curl -sL $GO_URL | tar -C /usr/local -xz
export PATH=/usr/local/go/bin:$PATH
echo 'export PATH=/usr/local/go/bin:$PATH' >/etc/profile.d/go.sh
GO_BIN=/usr/local/go/bin/go

# ── 系统依赖 ──────────────────────────────────────────────
apt update -qq
apt install -y --no-install-recommends git jq curl openssl

# ── 拉源码 ───────────────────────────────────────────────
REPO="https://github.com/hiatunnel/hia-tunnel"
INSTALL_DIR="/opt/hia-tunnel"
BIN_DIR="/usr/local/bin"
CONF_DIR="/etc/swift-tunnel"
rm -rf "$INSTALL_DIR"
git clone --depth 1 "$REPO" "$INSTALL_DIR"
cd "$INSTALL_DIR"

# 生成自签 TLS 证书供 embed
cd internal/transport
openssl req -x509 -newkey rsa:2048 -days 3650 -nodes \
  -subj "/CN=HiaTunnel Dev" \
  -keyout devkey.pem -out devcert.pem >/dev/null 2>&1
cd ../../

$GO_BIN mod download
$GO_BIN build -o hia-tunnel-server ./cmd/server

install -Dm755 hia-tunnel-server "$BIN_DIR/hia-tunnel-server"
install -Dm755 scripts/menu.sh     /usr/local/bin/hia-menu

# ── 默认配置 ─────────────────────────────────────────────
mkdir -p "$CONF_DIR"
cat > "$CONF_DIR/server.json" <<EOF
{
  "listen": ":${PORT}",
  "psk": "${PSK}",
  "forwards": []
}
EOF

install -Dm644 systemd/hia-tunnel-server.service \
               /etc/systemd/system/hia-tunnel-server.service
systemctl daemon-reload
systemctl enable --now hia-tunnel-server

echo -e "\n✅ 安装完成！使用  hia-menu  添加端口转发或修改配置。"
