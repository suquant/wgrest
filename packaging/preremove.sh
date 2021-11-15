#!/usr/bin/env bash
set -eu

systemctl stop wgrest.service || true
systemctl disable wgrest.service || true
