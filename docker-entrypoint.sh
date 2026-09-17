#!/bin/bash
set -uo pipefail

./autossl &
api_pid=$!

cd /root/web
./node_modules/.bin/react-router-serve ./build/server/index.js &
web_pid=$!

shutdown() {
  trap - TERM INT
  kill "$api_pid" "$web_pid" 2>/dev/null || true
  wait "$api_pid" "$web_pid" 2>/dev/null || true
}

trap shutdown TERM INT
wait -n "$api_pid" "$web_pid"
status=$?
shutdown
exit "$status"
