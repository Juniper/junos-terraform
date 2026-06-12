#!/bin/bash

set -e
set -x

jtaf-yang2go -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang -t vqfx-evpn-vxlan
