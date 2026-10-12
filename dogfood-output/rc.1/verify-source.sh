#!/usr/bin/env bash
set -euo pipefail
exec > /src/dogfood-output/rc.1/source-build.log 2>&1
task_dir=$(mktemp -d /tmp/meu-ubuntu-rc-source-XXXXXX)
tar -xzf /src/dist/meu-ubuntu-1.0.0-rc.1-source.tar.gz -C "$task_dir"
cd "$task_dir/meu-ubuntu-1.0.0-rc.1-source"
test ! -d private
test ! -d frontend/node_modules
python3 build/scripts/generate-docs.py --check
(cd frontend && npm ci --no-audit --no-fund && npm run build)
go build -trimpath -buildvcs=false -tags desktop,production,webkit2_41 -ldflags '-s -w -X main.version=1.0.0-rc.1' -o "$task_dir/meu-ubuntu" .
test "$("$task_dir/meu-ubuntu" --version)" = 'meu-ubuntu 1.0.0-rc.1'
"$task_dir/meu-ubuntu" --licenses | grep -F 'GNU GENERAL PUBLIC LICENSE' >/dev/null
printf 'Candidate successfully rebuilt from extracted source archive.\n'
chown 1000:1000 /src/dogfood-output/rc.1/source-build.log
