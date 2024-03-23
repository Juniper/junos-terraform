#!/bin/bash -eu
set -o pipefail

if [ -z "$PNAME" ]
then
  echo "PNAME not found in environment"
fi

if ! gpg --fingerprint --with-colons "$PNAME" 2>/dev/null | awk -F: '/fpr:/{print $10}' | head -1
then
  dirname=$(dirname -- "$(readlink -f -- "$0")")
  "${dirname}/create_key.sh" "$PNAME"
fi
