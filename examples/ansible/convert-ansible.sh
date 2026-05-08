#!/bin/bash

set -e

RESET_PAYLOAD_CACHE="${RESET_PAYLOAD_CACHE:-1}"


if [[ "$RESET_PAYLOAD_CACHE" == "1" ]]; then
	rm -f ansible_files/.jtaf_host_payloads.json
	rm -f ansible_files/group_vars/all.yaml
fi

jtaf-xml2yaml -x ../evpn-vxlan-dc/dc1/dc1-*leaf* ../evpn-vxlan-dc/dc1/dc1-*spine* ../evpn-vxlan-dc/dc2/dc2-*spine* -j ansible-provider-junos-vqfx-ansible-role/trimmed_schema.json -d ansible_files --grouping-hosts-file switches_grouping_hosts

jtaf-xml2yaml -x ../evpn-vxlan-dc/dc1/dc1-*firewall* ../evpn-vxlan-dc/dc2/dc2-*firewall* -j ansible-provider-junos-srx-ansible-role/trimmed_schema.json -d ansible_files --grouping-hosts-file firewall_grouping_hosts


