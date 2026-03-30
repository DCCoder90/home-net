#!/usr/bin/env bash
# Build all three custom Pulumi provider binaries and install them to ~/.pulumi/plugins/
# Run this script once before `pulumi up`, and again whenever a provider changes.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROVIDERS_DIR="$SCRIPT_DIR/../providers"
PLUGINS_DIR="${HOME}/.pulumi/plugins"

build_provider() {
  local name="$1"
  local version="${2:-0.0.1}"
  local dir="$PROVIDERS_DIR/$name"
  local dest="$PLUGINS_DIR/resource-${name}-v${version}"

  echo "==> Building provider: $name v$version"
  (
    cd "$dir"
    go mod tidy
    go build -o "pulumi-resource-${name}" .
  )

  echo "==> Installing $name to $dest"
  mkdir -p "$dest"
  cp "$dir/pulumi-resource-${name}" "$dest/"
  echo "    Installed: $dest/pulumi-resource-${name}"
}

mkdir -p "$PLUGINS_DIR"

build_provider "npmproxy"
build_provider "technitium"
build_provider "authentik"

echo ""
echo "All providers built and installed successfully."
echo "Installed plugins:"
ls "$PLUGINS_DIR" | grep -E "^resource-(npmproxy|technitium|authentik)"
