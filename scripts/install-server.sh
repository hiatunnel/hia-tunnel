#!/usr/bin/env bash
set -euo pipefail

### 固定默认监听端口，可通过环境变量 PORT 覆盖
PORT="${PORT:-8443}"

### 预共享密钥：如未提前导出则随机生成
PSK="${PSK:-$(openssl rand -hex 32)}"

echo -e "\n[Hia‑Tunnel] 监听端口 : ${PORT}"
echo -e "[Hia‑Tunnel] 预共享密钥 : ${PSK}\n"

# 1) 依赖
apt update -qq
apt install -y --no-install-recommends git golang-go jq curl openssl

# 2) 拉源码
REPO="https://github.com/hiatunnel/hia-tunnel"
INSTALL_DIR="/opt/hia-tunnel"
BIN_DIR="/usr/local/bin"
CONF_DIR="/etc/swift-tunnel"
rm -rf "$INSTALL_DIR"
git clone --depth 1 "$REPO" "$INSTALL_DIR"

# 3) 生成自签 TLS (embed)
cd "$INSTALL_DIR/internal/transport"
openssl req -x509 -newkey rsa:2048 -days 3650 -nodes \
  -subj "/CN=HiaTunnel Dev" \
  -keyout devkey.pem -out devcert.pem >/dev/null 2>&1
cd "$INSTALL_DIR"

# 4) 构建
go mod tidy -e
go build -o hia-tunnel-server ./cmd/server

install -Dm755 hia-tunnel-server "$BIN_DIR/hia-tunnel-server"
install -Dm755 scripts/menu.sh     /usr/local/bin/hia-menu

# 5) 默认配置（直接写 8443，可后续修改）
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

echo -e "\n✅ 安装完成！"
echo "   • 服务已在 ${PORT}/udp 监听"
echo "   • 运行  hia-menu  添加端口转发或修改监听端口 / PSK"
