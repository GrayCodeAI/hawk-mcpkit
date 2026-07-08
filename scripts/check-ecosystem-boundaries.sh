#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# hawk-mcpkit is a foundation library: it sits below every engine and below
# hawk itself. It must depend on nothing in the hawk ecosystem — only the
# Go standard library and the upstream mark3labs/mcp-go package.
FORBIDDEN='github\.com/GrayCodeAI/'

if command -v rg >/dev/null 2>&1; then
  violations="$(rg -n "$FORBIDDEN" --glob '*.go' . || true)"
else
  violations="$(grep -rn --include='*.go' -E "$FORBIDDEN" . || true)"
fi

if [[ -n "${violations}" ]]; then
  echo "forbidden hawk-eco imports found in hawk-mcpkit:"
  echo "${violations}"
  echo
  echo "hawk-mcpkit is a foundation repo — it must not depend on hawk, engines, or any other GrayCodeAI/* package"
  exit 1
fi

echo "ecosystem boundary guard passed"
