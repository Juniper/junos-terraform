#!/bin/bash

set -e

jtaf-xml2yaml -x ../evpn-vxlan-dc/dc1/dc1-*leaf* ../evpn-vxlan-dc/dc1/dc1-*spine* ../evpn-vxlan-dc/dc2/dc2-*spine*  -j ansible-provider-junos-vqfx-ansible-role/trimmed_schema.json -d vqfx_ansible_files

jtaf-xml2yaml -x ../evpn-vxlan-dc/dc1/dc1-*firewall* ../evpn-vxlan-dc/dc2/dc2-*firewall*  -j ansible-provider-junos-srx-ansible-role/trimmed_schema.json -d srx_ansible_files
