#!/usr/bin/env bash
set -e
systemctl disable --now hia-tunnel-server hia-tunnel-client || true
rm -f /usr/local/bin/hia-tunnel-{server,client} /usr/local/bin/hia-*
rm -rf /opt/hia-tunnel /etc/swift-tunnel /etc/systemd/system/hia-tunnel-*.service
systemctl daemon-reload
echo "Hia Tunnel removed."
