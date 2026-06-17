#!/bin/bash

# One-time setup script for ActivityPub federation tests
# Run with: sudo ./setup.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# When a GitHub token is available (e.g. GITHUB_TOKEN forwarded into this
# script in CI), authenticate downloads from github.com so they use the
# 5,000/hour rate limit instead of the anonymous, shared-per-IP 60/hour limit
# that frequently fails on CI runners. Local runs without a token download
# anonymously. curl drops the Authorization header on the cross-host redirect
# to the asset CDN, so the token is not leaked.
gh_curl_auth=()
if [[ -n "${GITHUB_TOKEN:-}" ]]; then
    gh_curl_auth=(-H "Authorization: Bearer ${GITHUB_TOKEN}")
fi

if [[ ${EUID} -ne 0 ]]; then
    log_error "This script must be run as root (sudo ./setup.sh)"
    exit 1
fi

# Get the actual user who ran sudo
ACTUAL_USER="${SUDO_USER:-${USER}}"
ACTUAL_HOME=$(getent passwd "${ACTUAL_USER}" | cut -d: -f6)

log_info "Setting up ActivityPub federation test environment..."

# 1. Install mkcert
if ! command -v mkcert &> /dev/null; then
    log_info "Downloading mkcert..."
    curl -sL "${gh_curl_auth[@]}" "https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-linux-amd64" -o /usr/local/bin/mkcert
    chmod +x /usr/local/bin/mkcert
else
    log_info "mkcert already installed"
fi

# 2. Install CA into system trust store
log_info "Installing mkcert CA into system trust store..."
sudo -u "${ACTUAL_USER}" CAROOT="${ACTUAL_HOME}/.local/share/mkcert" mkcert -install

# 3. Generate certificates
log_info "Generating certificates..."
mkdir -p "${SCRIPT_DIR}/certs"
chown "${ACTUAL_USER}:${ACTUAL_USER}" "${SCRIPT_DIR}/certs"
sudo -u "${ACTUAL_USER}" CAROOT="${ACTUAL_HOME}/.local/share/mkcert" mkcert \
    -cert-file "${SCRIPT_DIR}/certs/cert.pem" \
    -key-file "${SCRIPT_DIR}/certs/key.pem" \
    owncast.local owncast2.local snac.local localhost 127.0.0.1

# 4. Add hosts entries
if ! grep -q "owncast.local" /etc/hosts; then
    log_info "Adding hosts entries..."
    echo "127.0.0.1 owncast.local owncast2.local snac.local" >> /etc/hosts
elif ! grep -q "owncast2.local" /etc/hosts; then
    log_info "Adding owncast2.local hosts entry..."
    echo "127.0.0.1 owncast2.local" >> /etc/hosts
else
    log_info "Hosts entries already exist"
fi

# 5. Check for snac2
if ! command -v snac &> /dev/null; then
    log_error "snac2 is not installed. Please install it manually."
    log_error "See: https://codeberg.org/grunfink/snac2"
    exit 1
else
    log_info "snac2 found: $(command -v snac)"
fi

# 6. Install Caddy
if ! command -v caddy &> /dev/null; then
    log_info "Installing Caddy..."
    CADDY_VERSION="2.8.4"
    curl -sL "${gh_curl_auth[@]}" "https://github.com/caddyserver/caddy/releases/download/v${CADDY_VERSION}/caddy_${CADDY_VERSION}_linux_amd64.tar.gz" | tar -xz -C /usr/local/bin caddy
    chmod +x /usr/local/bin/caddy
    log_info "Caddy installed"
else
    log_info "Caddy already installed: $(command -v caddy)"
fi

echo ""
log_info "Setup complete!"
echo ""
echo "You can now run the test:"
echo "  cd ${SCRIPT_DIR}"
echo "  ./run.sh"
