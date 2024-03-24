#!/bin/bash -eu
set -o pipefail

DIRNAME=$(dirname -- "$(readlink -f -- "$0")")
DIST_DIR="$DIRNAME/../dist"

export PNAME=$("${DIRNAME}/provider_name.sh")
export GPG_FINGERPRINT=$(${DIRNAME}/key_fingerprint.sh $PNAME)

go run github.com/goreleaser/goreleaser release --clean --snapshot

gpg --armor --export "$GPG_FINGERPRINT" > "$DIST_DIR/gpg_key.asc"
echo "$GPG_FINGERPRINT" > "$DIST_DIR/gpg_key.id"
