#!/usr/bin/env bash
set -euo pipefail
cd /src
# Pin the checker; its vulnerability database is refreshed on every audit.
go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 -tags webkit2_41 ./...
(cd frontend && npm audit --audit-level=high)
