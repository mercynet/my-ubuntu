#!/usr/bin/env bash
# Optional independent audit. Requires an authenticated Claude Code CLI.
set -euo pipefail
cd "$(dirname "$0")/../.."
MODEL="${REVIEW_MODEL:-claude-opus-5-5}"
OUTPUT="${REVIEW_OUTPUT:-build/test/out/adversarial-review.json}"
mkdir -p "$(dirname "$OUTPUT")"

claude -p --model "$MODEL" --effort high \
  --restricted --safe-mode --permission-mode dontAsk --permission-prompts none \
  --allowedTools Read,Glob,Grep --tools Read,Glob,Grep \
  --strict-mcp-config --mcp-config '{"mcpServers":{}}' \
  --disable-slash-commands --no-session-persistence --output-format json \
  < build/review/adversarial-prompt.txt > "$OUTPUT"

python3 - "$OUTPUT" <<'PY'
import json, sys
with open(sys.argv[1]) as source:
    result = json.load(source)
print(result.get("result", "Sem resposta textual."))
if result.get("is_error"):
    raise SystemExit(1)
PY
