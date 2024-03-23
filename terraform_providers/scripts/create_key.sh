#!/bin/bash -eu
set -o pipefail

if [ -z "$PNAME" ]
then
  echo "PNAME not found in environment"
fi

gpg --batch --passphrase '' --quick-gen-key "$PNAME" rsa sign > /dev/null 2>&1

gpg --fingerprint --with-colons "$PNAME" 2>/dev/null | awk -F: '/fpr:/{print $10}' | head -1
