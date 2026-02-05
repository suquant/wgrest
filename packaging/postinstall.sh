#!/usr/bin/env bash
set -eu

# Create data directory for certs cache
mkdir -p /var/lib/wgrest
chmod 700 /var/lib/wgrest

# Enable and start service
systemctl daemon-reload
systemctl enable wgrest.service
