#!/bin/bash
# Build providers whose embedded schema is trimmed to the example XML.
# Same provider source as build-generic.sh; only the schema scope differs.
#
# Usage: cd examples/providers && bash build.sh

set -e

providers_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "$providers_dir/../.." && pwd)"
output_dir="${JTAF_PROVIDER_OUTPUT_DIR:-$providers_dir}"
mkdir -p "$output_dir"
cd "$output_dir"

jtaf-yang2go \
        -p "$repo_root/examples/yang/18.2/18.2R3/common" \
        "$repo_root"/examples/yang/18.2/18.2R3/junos-qfx/conf/*.yang \
        -x "$repo_root"/examples/evpn-vxlan-dc/dc1/dc1-*leaf* \
                 "$repo_root"/examples/evpn-vxlan-dc/dc1/dc1-*spine* \
                 "$repo_root"/examples/evpn-vxlan-dc/dc2/dc2-*spine* \
        -t vqfx-evpn-vxlan

jtaf-yang2go \
        -p "$repo_root/examples/yang/18.2/18.2R3/common" \
        "$repo_root"/examples/yang/18.2/18.2R3/junos-es/conf/*.yang \
        -x "$repo_root"/examples/evpn-vxlan-dc/dc1/dc1-*firewall* \
                 "$repo_root"/examples/evpn-vxlan-dc/dc2/dc2-*firewall* \
        -t vsrx-evpn-vxlan

for provider_type in vqfx-evpn-vxlan vsrx-evpn-vxlan; do
        echo "Compiling $provider_type provider..."
        (cd "terraform-provider-junos-$provider_type" && go build .)
done
