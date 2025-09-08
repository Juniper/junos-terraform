#!/bin/bash

set -x
set -e

(cd ../../.. && git clone https://github.com/Juniper/yang.git 2>&1 | grep "destination path 'yang' already exists")

QFXFILES=`echo ../evpn-vxlan-dc/dc1/dc1-*leaf* ../evpn-vxlan-dc/dc1/dc1-*spine* ../evpn-vxlan-dc/dc2/dc2-*spine*`
QFXXMLARGS=$(
	for i in $QFXFILES; do
		echo "-x $i"
	done
)

jtaf-yang2go -p ../../../yang/18.2/18.2R3/common ../../../yang/18.2/18.2R3/junos-qfx/conf/*.yang $QFXXMLARGS -t qfx-evpn-vxlan 

# this is really not great
mv trimmed_schema.json terraform-provider-junos-qfx-evpn-vxlan

SRXFILES=`echo ../evpn-vxlan-dc/dc1/dc1-*firewall* ../evpn-vxlan-dc/dc2/dc2-*firewall*`
SRXXMLARGS=$(
	for i in $SRXFILES; do
		echo "-x $i"
	done
)

jtaf-yang2go -p ../../../yang/18.2/18.2R3/common ../../../yang/18.2/18.2R3/junos-es/conf/*.yang $SRXXMLARGS -t srx-evpn-vxlan 

# this is really not great
mv trimmed_schema.json terraform-provider-junos-srx-evpn-vxlan
