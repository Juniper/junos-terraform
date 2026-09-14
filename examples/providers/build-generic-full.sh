#!/bin/bash
# Build a generic provider from the full YANG model (no -x filtering).
# This creates a single universal provider covering every QFX YANG node.
#
# Usage: cd examples/providers && bash build-generic-full.sh

set -e

echo "Building full-model QFX generic provider (no -x filtering)..."
echo "  pyang will process all 33 QFX YANG files — this takes ~2 minutes."
echo ""

jtaf-yang2go --generic \
  -p ../yang/18.2/18.2R3/common \
  ../yang/18.2/18.2R3/junos-qfx/conf/*.yang \
  -t qfx-full

echo ""
echo "Compiling full-model provider..."
cd terraform-provider-junos-qfx-full
go build .

SCHEMA_SIZE=$(du -sh trimmed_schema.json | cut -f1)
BINARY_SIZE=$(du -sh terraform-provider-junos-qfx-full | cut -f1)
echo ""
echo "  Schema:  $SCHEMA_SIZE"
echo "  Binary:  $BINARY_SIZE"

echo ""
echo "Installing to \$GOPATH/bin..."
go install .
cd ..

echo ""
echo "Done. Full-model provider installed to $(go env GOPATH)/bin/"
echo ""
echo "This provider covers every Junos QFX 18.2 configuration node."
echo "Only rebuild when a new Junos/YANG release is available."
