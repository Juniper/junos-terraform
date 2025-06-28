#!/bin/bash -eu
set -o pipefail

[ $# -eq 1 ]

if ! gpg --fingerprint --with-colons "$1" 2>/dev/null | awk -F: '/fpr:/{print $10}' | head -1
then
  DIRNAME=$(dirname -- "$(readlink -f -- "$0")")
  "${DIRNAME}/create_key.sh" "$1"
fi
