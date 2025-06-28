#!/bin/bash -eu
set -o pipefail

DIRNAME=$(dirname -- "$(readlink -f -- "$0")")
DIST_DIR="$DIRNAME/../dist"

export PTYPE=$("${DIRNAME}/provider_name.sh")
export GPG_FINGERPRINT=$(${DIRNAME}/key_fingerprint.sh $PTYPE)

go run github.com/goreleaser/goreleaser release --clean --snapshot

gpg --armor --export "$GPG_FINGERPRINT" > "$DIST_DIR/gpg_key.asc"
echo "$GPG_FINGERPRINT" > "$DIST_DIR/gpg_key.id"

go run github.com/chrismarget/lambda-tf-registry/cmd/register

"${DIRNAME}/delete_key.sh" "$GPG_FINGERPRINT"
