#!/bin/bash

set -e

providers_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "$providers_dir/../.." && pwd)"
provider_root="${JTAF_PROVIDER_OUTPUT_DIR:-$providers_dir}"
terraform_dir="${JTAF_TERRAFORM_OUTPUT_DIR:-$repo_root/examples/terraform_files}"
mkdir -p "$terraform_dir"

jtaf-xml2tf \
	-x "$repo_root"/examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml \
		 "$repo_root"/examples/evpn-vxlan-dc/dc2/*spine*.xml \
	-j "$provider_root/terraform-provider-junos-vqfx-evpn-vxlan/trimmed_schema.json.gz" \
	-t vqfx-evpn-vxlan \
	-d "$terraform_dir" \
	-u jcluser \
	-p 'Juniper!1'

#jtaf-xml2tf -x ../evpn-vxlan-dc/dc1/*firewall*.xml ../evpn-vxlan-dc/dc2/*firewall*.xml -j "$provider_root/terraform-provider-junos-vsrx-evpn-vxlan/trimmed_schema.json.gz" -t vsrx-evpn-vxlan -d "$terraform_dir" -u jcluser -p 'Juniper!1'

