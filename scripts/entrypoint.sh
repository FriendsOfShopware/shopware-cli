#!/usr/bin/env bash
set -euo pipefail

source /root/.bashrc
nvm use "${NODE_VERSION}" > /dev/null

exec "$@"
