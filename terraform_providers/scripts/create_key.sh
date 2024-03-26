#!/bin/bash -eu
set -o pipefail

[ $# -eq 1 ]

gpg --batch --passphrase '' --quick-gen-key "$1" rsa sign > /dev/null 2>&1
gpg --fingerprint --with-colons "$1" 2>/dev/null | awk -F: '/fpr:/{print $10}' | head -1
