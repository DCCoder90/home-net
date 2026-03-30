#!/usr/bin/env bash
# Generate `pulumi import` commands for resources currently managed by Terraform.
# Review the output before running — resource IDs must match what Terraform has in state.
#
# Prerequisites:
#   - terraform (tofu) CLI available
#   - `pulumi` CLI logged in
#   - Run from the repo root or adjust paths
#
# Usage:
#   ./scripts/import-state.sh > import.sh
#   chmod +x import.sh && ./import.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TF_CORE="$SCRIPT_DIR/../../terraform/core"
TF_APPS="$SCRIPT_DIR/../../terraform/apps"

echo "#!/usr/bin/env bash"
echo "# Auto-generated import commands. Review IDs before running."
echo "set -euo pipefail"
echo ""
echo "cd \"$SCRIPT_DIR/..\""
echo ""

# Helper to get TF output, silently ignoring failures.
tf_output() {
  local dir="$1" name="$2"
  (cd "$dir" && terraform output -raw "$name" 2>/dev/null) || echo ""
}

# Helper to get a list of TF resource IDs from state.
tf_state_show() {
  local dir="$1" addr="$2"
  (cd "$dir" && terraform state show "$addr" 2>/dev/null) || echo ""
}

# ---- Core: Technitium DNS container ----
echo "# Core: Technitium DNS container"
echo "pulumi import 'docker:index/container:Container' technitium_dns technitium_dns"
echo ""

# ---- Core: Nginx Proxy Manager container ----
echo "# Core: Nginx Proxy Manager container"
echo "pulumi import 'docker:index/container:Container' nginx_proxy_manager Nginx-Proxy-Manager-Official"
echo ""

# ---- NPM Access Lists ----
echo "# NPM Access Lists (get IDs from NPM UI or API)"
echo "# pulumi import 'npmproxy:index:AccessList' npm-acl-internal <internal-acl-id>"
echo "# pulumi import 'npmproxy:index:AccessList' npm-acl-cloudflare <cloudflare-acl-id>"
echo ""

# ---- Technitium DNS Records ----
# Each record ID in our provider is "domain|zone|type|ttl"
echo "# Technitium DNS Records (format: domain|zone|type|ttl)"
echo "# pulumi import 'technitium:index:Record' <service-name>-dns '<domain>|<zone>|A|3600'"
echo ""

# ---- NPM Certificates ----
echo "# NPM Certificates (get IDs from NPM UI or API)"
echo "# pulumi import 'npmproxy:index:Certificate' <service-name>-cert <cert-id>"
echo ""

# ---- NPM Proxy Hosts ----
echo "# NPM Proxy Hosts (get IDs from NPM UI or API)"
echo "# pulumi import 'npmproxy:index:ProxyHost' <service-name>-proxy <proxy-host-id>"
echo ""

# ---- Docker Containers (from stacks and services) ----
echo "# Docker Containers (service containers)"
echo "# pulumi import 'docker:index/container:Container' <service-name> <container-name>"
echo ""

# ---- Authentik resources ----
echo "# Authentik Groups (get PKs from Authentik UI)"
echo "# pulumi import 'authentik:index:Group' <service-name>-group <group-pk>"
echo ""
echo "# Authentik Applications (get UUIDs from Authentik UI)"
echo "# pulumi import 'authentik:index:Application' <service-name>-app <app-uuid>"
echo ""
echo "# Authentik Proxy Providers (get PKs from Authentik UI)"
echo "# pulumi import 'authentik:index:ProviderProxy' <service-name>-proxy-provider <provider-pk>"
echo ""
echo "# Authentik OAuth2 Providers (get PKs from Authentik UI)"
echo "# pulumi import 'authentik:index:ProviderOAuth2' <service-name>-oauth-provider <provider-pk>"
echo ""
echo "# Authentik Policy Bindings (get PKs from Authentik UI)"
echo "# pulumi import 'authentik:index:PolicyBinding' <service-name>-binding <binding-pk>"
echo ""

echo "echo 'Import complete. Run: pulumi up --refresh'"
