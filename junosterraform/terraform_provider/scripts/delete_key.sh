#!/bin/bash -eu
set -o pipefail

[ $# -eq 1 ]

gpg --delete-secret-and-public-key --batch --yes "$1"
