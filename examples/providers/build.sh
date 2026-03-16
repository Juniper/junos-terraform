#!/bin/bash

set -e

jtaf-yang2go -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang -x ../evpn-vxlan-dc/dc1/dc1-*leaf* ../evpn-vxlan-dc/dc1/dc1-*spine* ../evpn-vxlan-dc/dc2/dc2-*spine* -t vqfx-evpn-vxlan

jtaf-yang2go -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-es/conf/*.yang -x ../evpn-vxlan-dc/dc1/dc1-*firewall* ../evpn-vxlan-dc/dc2/dc2-*firewall* -t vsrx-evpn-vxlan
