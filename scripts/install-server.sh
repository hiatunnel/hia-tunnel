#!/usr/bin/env bash
set -euo pipefail

while [[ -z ${PORT:-} ]]; do
  read -rp "⚙️  请输入监听端口(必填, 如 8443): " PORT
done
PSK="${PSK:-$(openssl rand -hex 32)}"
echo -e "\n监听端口: $PORT"
echo -e "预共享密钥: $PSK\n"

apt update -qq
apt install -y --no-install-recommends git golang-go jq curl

REPO="https://github.com/hiatunnel/hia-tunnel"
INSTALL_DIR="/opt/hia-tunnel"
BIN_DIR="/usr/local/bin"
CONF_DIR="/etc/swift-tunnel"

rm -rf "$INSTALL_DIR"
git clone --depth 1 "$REPO" "$INSTALL_DIR"
cd "$INSTALL_DIR"
go build -o hia-tunnel-server ./cmd/server
install -Dm755 hia-tunnel-server "$BIN_DIR/hia-tunnel-server"
install -Dm755 scripts/menu.sh /usr/local/bin/hia-menu

mkdir -p "$CONF_DIR"
cat > "$CONF_DIR/server.json" <<EOF
{
  "listen": ":${PORT}",
  "psk": "${PSK}",
  "forwards": []
}
EOF

install -Dm644 systemd/hia-tunnel-server.service /etc/systemd/system/hia-tunnel-server.service
systemctl daemon-reload
systemctl enable --now hia-tunnel-server

echo -e "\n✅ 安装完成，使用 'hia-menu' 管理转发与 PSK。"
