#!/usr/bin/env bash
set -eu

function remove_user() {
  deluser --quiet --system wgrest >/dev/null ||
    echo "Failed to remove user"
}

case $@ in
# apt purge passes "purge"
"purge")
  remove_user
  ;;
  # apt remove passes "remove"
"remove") ;;

esac
