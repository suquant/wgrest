#!/usr/bin/env bash
set -eu

adduser --system wgrest --home /var/lib/wgrest

systemctl enable "/etc/systemd/system/wgrest.service"
