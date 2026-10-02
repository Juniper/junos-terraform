#!/bin/bash
# Convert the group-based example into Terraform configuration, keeping the
# group hierarchy rather than flattening what the devices inherit.
#
# Usage: cd examples/providers && bash convert-groups.sh

set -e

providers_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "$providers_dir/../.." && pwd)"
provider_root="${JTAF_PROVIDER_OUTPUT_DIR:-$providers_dir}"
terraform_dir="${JTAF_TERRAFORM_OUTPUT_DIR:-$repo_root/examples/terraform_files_groups}"
mkdir -p "$terraform_dir"

jtaf-xml2tf \
        --groups --generic \
	-x "$repo_root"/examples/evpn-vxlan-dc-groups/dc1/*{spine,leaf}*.xml \
		 "$repo_root"/examples/evpn-vxlan-dc-groups/dc2/*spine*.xml \
        -j "$provider_root/terraform-provider-junos-vqfx-evpn-vxlan-groups/trimmed_schema.json.gz" \
        -t vqfx-evpn-vxlan-groups \
        -d "$terraform_dir" \
        -u jcluser \
        -p 'Juniper!1'
