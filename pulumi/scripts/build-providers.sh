#!/usr/bin/env bash
# build-providers.sh — compiles the custom Pulumi provider binaries and places
# them in pulumi/bin/ where Pulumi.yaml's plugins block can find them.
#
# Usage: run from the repo root or the pulumi/ directory.
#   ./pulumi/scripts/build-providers.sh
#
# Requirements: Go 1.22+ must be on PATH.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PULUMI_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_DIR="$PULUMI_DIR/bin"

mkdir -p "$BIN_DIR"

build_provider() {
  local name="$1"           # e.g. "technitium"
  local src="$PULUMI_DIR/providers/$name"

  echo "==> Building pulumi-resource-$name"
  pushd "$src" > /dev/null

  # Resolve dependencies on first build.
  go mod tidy

  # Build the provider binary for Linux (Pulumi runs in WSL).
  GOOS=linux GOARCH=amd64 \
    go build -o "$BIN_DIR/pulumi-resource-$name" .

  popd > /dev/null
  echo "    -> $BIN_DIR/pulumi-resource-$name"
}

build_provider technitium
build_provider npmproxy

# Also run go mod tidy for the main pulumi module so all dependencies are locked.
echo "==> Tidying main pulumi module"
pushd "$PULUMI_DIR" > /dev/null
go mod tidy
popd > /dev/null

echo ""
echo "All providers built successfully."
echo "You can now run: cd pulumi && pulumi up"
