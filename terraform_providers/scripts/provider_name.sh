#!/bin/bash -eu
set -o pipefail

dirname=$(dirname -- "$(readlink -f -- "$0")")
xpath_file="${dirname}/../../xpath_inputs.xml"

sum=$(shasum -a 256 "$xpath_file" | cut -b-8)
epoch=$(printf '%x' $(date +%s))

echo "${epoch}${sum}"
