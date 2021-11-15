#!/usr/bin/env bash
set -eu

if systemctl status wgrest &> /dev/null; then
    systemctl stop wgrest.service
    systemctl disable wgrest.service
fi
