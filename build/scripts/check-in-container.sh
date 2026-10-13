#!/usr/bin/env bash
# The same checks run locally and in CI, including a fresh checkout without assets.
set -euo pipefail
cd /src

restore_owner() {
  if [ -n "${HOST_UID:-}" ]; then
    chown -R "${HOST_UID}:${HOST_GID:-$HOST_UID}" frontend/dist frontend/node_modules 2>/dev/null || true
  fi
}
trap restore_owner EXIT

# Wails embeds frontend/dist, so compile it before checking the complete Go app.
python3 build/scripts/check-version.py
python3 build/scripts/generate-docs.py --check
python3 -m unittest discover -s build/tests -p 'test_*.py'
(cd frontend && npm ci --no-audit --no-fund && npm test && npm run build)
go vet -tags webkit2_41 ./...
go test -race -tags webkit2_41 ./...
