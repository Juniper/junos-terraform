#!/bin/bash -eu
set -o pipefail

[ $# -eq 1 ]

if ! gpg --fingerprint --with-colons "$1" 2>/dev/null | awk -F: '/fpr:/{print $10}' | head -1
then
  dirname=$(dirname -- "$(readlink -f -- "$0")")
  "${dirname}/create_key.sh" "$1"
fi