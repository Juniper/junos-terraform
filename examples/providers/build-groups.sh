#!/bin/bash
# Build providers that manage Junos configuration groups, with the embedded
# schema trimmed to the group-based example.
#
# Usage: cd examples/providers && bash build-groups.sh

set -e

providers_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "$providers_dir/../.." && pwd)"
output_dir="${JTAF_PROVIDER_OUTPUT_DIR:-$providers_dir}"
mkdir -p "$output_dir"
cd "$output_dir"

echo "Building groups-aware QFX provider..."
jtaf-yang2go --groups --generic \
        -p "$repo_root/examples/yang/18.2/18.2R3/common" \
        "$repo_root"/examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang \
        -t vqfx-evpn-vxlan-groups

echo "Compiling groups-aware QFX provider..."
(cd terraform-provider-junos-vqfx-evpn-vxlan-groups && go build .)

if [[ "${JTAF_SKIP_PROVIDER_INSTALL:-0}" == "1" ]]; then
        echo "Skipping global provider installation."
else
        (cd terraform-provider-junos-vqfx-evpn-vxlan-groups && go install .)
        echo "Installed to $(go env GOPATH)/bin/"
fi

echo
echo "Next: bash convert-groups.sh"
