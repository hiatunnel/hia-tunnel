#!/usr/bin/env bash
set -e
CONF="${1:-./client.json}"
TMP=$(mktemp)

list() { jq -r '.peers[]? | "☆ " + .name + " -> " + .server' "$CONF"; }

echo "=== Hia Tunnel 客户端菜单 ==="
PS3="选择操作: "
select opt in "添加/修改对接" "删除对接" "列出对接" "重启客户端" "退出"; do
  case $REPLY in
    1)
      read -rp "Name: " n
      read -rp "Server(host:port): " s
      read -rp "PSK(64hex): " k
      read -rp "本地 SOCKS 端口(如 1080): " p
      jq --arg n "$n" --arg s "$s" --arg k "$k" --arg p "127.0.0.1:$p" '
        .peers |= (map(select(.name==$n))|length==0
          ? (. + [{"name":$n,"server":$s,"psk":$k,"socks_local":$p,"max_streams":1024}])
          : (map(if .name==$n then .server=$s|.psk=$k|.socks_local=$p else . end)))'         "$CONF" >"$TMP" && mv "$TMP" "$CONF"
      ;;
    2)
      list; read -rp "要删除的 Name: " dn
      jq --arg dn "$dn" 'del(.peers[]|select(.name==$dn))' "$CONF" >"$TMP" && mv "$TMP" "$CONF"
      ;;
    3) list ;;
    4) systemctl restart hia-tunnel-client || echo "请手动重启" ;;
    *) ;;
  esac
  break
done
