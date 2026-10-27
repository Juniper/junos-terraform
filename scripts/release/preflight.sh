#!/bin/bash
# Fail fast before a release build: the tools it needs, free space on the
# filesystem that holds the build directory, and (as a warning) memory.
#
# Environment: BUILD_DIR MIN_FREE_GB JOBS AUTO_CLEAN REQUIRED_TOOLS GB_PER_JOB
set -e -o pipefail

BUILD_DIR=${BUILD_DIR:-build/release}
MIN_FREE_GB=${MIN_FREE_GB:-5}
JOBS=${JOBS:-2}
AUTO_CLEAN=${AUTO_CLEAN:-0}
# Generating one full-model provider peaks at 6-8 GB (pyang, then the schema
# compiler), measured on the 25.4R1 QFX and EVO models.
GB_PER_JOB=${GB_PER_JOB:-8}
REQUIRED_TOOLS=${REQUIRED_TOOLS:-git pyang jtaf-yang2go jtaf-provider go tar python3}
here=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
CLEAN_SPACE_SCRIPT=${CLEAN_SPACE_SCRIPT:-$here/clean_space.sh}

fail() { echo "preflight: $*" >&2; exit 1; }

# ---- tools ---------------------------------------------------------------
missing=()
for tool in $REQUIRED_TOOLS; do
  command -v "$tool" >/dev/null 2>&1 || missing+=("$tool")
done
if ((${#missing[@]})); then
  echo "preflight: missing tool(s): ${missing[*]}" >&2
  fail "pyang and the jtaf-* commands come from 'pip install -e .' in this repository; go and git from your package manager"
fi

# ---- disk ----------------------------------------------------------------
mkdir -p "$BUILD_DIR/tmp"
free_gb() { df -Pk "$1" | awk 'NR==2 {printf "%d", $4 / 1024 / 1024}'; }

free=$(free_gb "$BUILD_DIR")
if (( free < MIN_FREE_GB )) && [[ "$AUTO_CLEAN" == 1 ]]; then
  echo "preflight: ${free} GB free under $BUILD_DIR, need ${MIN_FREE_GB} GB; AUTO_CLEAN=1, running clean-space"
  BUILD_DIR=$BUILD_DIR bash "$CLEAN_SPACE_SCRIPT"
  free=$(free_gb "$BUILD_DIR")
fi
if (( free < MIN_FREE_GB )); then
  echo "preflight: ${free} GB free on the filesystem holding $BUILD_DIR, need at least ${MIN_FREE_GB} GB (MIN_FREE_GB)" >&2
  fail "reclaim space with 'make clean-space' (DRY_RUN=1 to preview), or rerun with AUTO_CLEAN=1"
fi

# ---- memory (warning only) -----------------------------------------------
if [[ -r /proc/meminfo ]]; then
  avail=$(awk '/^MemAvailable:/ {printf "%d", $2 / 1024 / 1024}' /proc/meminfo)
  need=$(( JOBS * GB_PER_JOB ))
  if (( avail < need )); then
    echo "preflight: warning: ${avail} GB of memory available; generating ${JOBS} providers at once peaks around ${need} GB. Consider JOBS=1." >&2
  fi
fi

echo "preflight: ok (${free} GB free under $BUILD_DIR; tools: $REQUIRED_TOOLS)"
