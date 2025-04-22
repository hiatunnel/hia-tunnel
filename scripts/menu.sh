#!/usr/bin/env bash
set -e
CONF="/etc/swift-tunnel/server.json"
TMP=$(mktemp)

list() { jq -r '.forwards[]? | "☆ " + .listen + " -> " + .target' "$CONF"; }

echo "=== Hia Tunnel 服务器菜单 ==="
PS3="选择操作: "
select opt in "修改监听端口/PSK" "添加转发" "删除转发" "列出转发" "重启服务" "退出"; do
  case $REPLY in
    1)
      port=""
      while [[ -z $port ]]; do read -rp "新监听端口: " port; done
      read -rp "新 PSK(64hex 留空不变): " psk
      jq --arg p ":$port" --arg k "$psk" '
        .listen=$p | if $k != "" then .psk=$k else . end' "$CONF" >"$TMP" && mv "$TMP" "$CONF"
      ;;
    2)
      read -rp "监听(0.0.0.0:8022): " l
      read -rp "目标(192.168.1.50:22): " t
      jq --arg l "$l" --arg t "$t" '.forwards += [{"listen":$l,"target":$t}]' "$CONF" >"$TMP" && mv "$TMP" "$CONF"
      ;;
    3)
      list
      read -rp "删除的 listen: " d
      jq --arg d "$d" 'del(.forwards[]|select(.listen==$d))' "$CONF" >"$TMP" && mv "$TMP" "$CONF"
      ;;
    4) list ;;
    5) systemctl restart hia-tunnel-server ;;
    *) break ;;
  esac
  break
done
