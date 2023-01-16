#!/usr/bin/env bash
set -euo pipefail

source /root/.bashrc
nvm use "${NODE_VERSION:-18}" > /dev/null

exec "$@"
