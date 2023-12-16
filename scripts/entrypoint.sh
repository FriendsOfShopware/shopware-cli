#!/usr/bin/env bash
set -euo pipefail

source /root/.bashrc

if ! find /root/.nvm/versions/node/ -maxdepth 1 -name "v${NODE_VERSION}*" | grep . &>/dev/null; then
  nvm install "${NODE_VERSION}" --silent &>/dev/null
fi

nvm use "${NODE_VERSION}" --silent &>/dev/null

exec "$@"
