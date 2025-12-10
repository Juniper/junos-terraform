#!/bin/bash

set -e

jtaf-xml2yaml -x ../ansible/*interfaces.xml -j ansible-provider-junos-qfx-ansible-test/trimmed_schema.json -d ansible_files

