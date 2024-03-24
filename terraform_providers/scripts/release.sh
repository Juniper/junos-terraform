#!/bin/bash -eu
set -o pipefail

DIRNAME=$(dirname -- "$(readlink -f -- "$0")")

export PNAME=$("${DIRNAME}/provider_name.sh")
export GPG_FINGERPRINT=$(${DIRNAME}/key_fingerprint.sh $PNAME)

go run github.com/goreleaser/goreleaser release --clean --snapshot
