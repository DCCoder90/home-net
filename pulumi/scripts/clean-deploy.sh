#!/usr/bin/env bash
# clean-deploy.sh — wipes provider binaries, does a clean rebuild, and runs pulumi up.
# Ensures Pulumi uses freshly compiled providers rather than any stale artifacts.
#
# Usage: run from the repo root or the pulumi/ directory.
#   ./pulumi/scripts/clean-deploy.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PULUMI_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_DIR="$PULUMI_DIR/bin"

echo "==> Cleaning provider binaries"
rm -f "$BIN_DIR"/pulumi-resource-*
echo "    -> $BIN_DIR cleaned"

echo ""
echo "==> Rebuilding providers"
bash "$SCRIPT_DIR/build-providers.sh"

echo ""
echo "==> Running pulumi up"
cd "$PULUMI_DIR"
pulumi up
