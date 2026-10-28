#!/bin/bash
# Reclaim space on the build machine without disturbing anything that is
# running. By default: this repository's stale release build output, stopped
# containers that exited more than PRUNE_UNTIL ago, dangling images, and
# unreferenced image build cache. Anything more is opt-in.
#
#   DRY_RUN=1                 list what would be removed; remove nothing
#   BUILD_DIR=build/release   this repository's release build directory
#   KEEP_RELEASES="25.4R1 25.4R1-EVO"   release build directories to keep
#   PRUNE_UNTIL=24h           only containers/build cache older than this
#   CLEAN_EXAMPLE_PROVIDERS=1 also the example scripts' provider output
#   PRUNE_UNUSED_IMAGES=1     also tagged images no container uses
#     PRUNE_IMAGE_AGE=168h      ...only those older than this (empty: all)
#   CLEAN_GO_CACHE=1          also `go clean -cache`
#   CLEAN_GO_MODCACHE=1       also `go clean -modcache` (asks unless FORCE=1)
#   OWN_TMP=1                 also this user's /tmp entries older than 7 days
#
# Never: docker system prune, volume prune, or stop/kill/rm of a running
# container.
set -e -o pipefail

DRY_RUN=${DRY_RUN:-0}
BUILD_DIR=${BUILD_DIR:-build/release}
KEEP_RELEASES=${KEEP_RELEASES:-}
PRUNE_UNTIL=${PRUNE_UNTIL:-24h}
CLEAN_EXAMPLE_PROVIDERS=${CLEAN_EXAMPLE_PROVIDERS:-0}
PRUNE_UNUSED_IMAGES=${PRUNE_UNUSED_IMAGES:-0}
PRUNE_IMAGE_AGE=${PRUNE_IMAGE_AGE-168h}
CLEAN_GO_CACHE=${CLEAN_GO_CACHE:-0}
CLEAN_GO_MODCACHE=${CLEAN_GO_MODCACHE:-0}
OWN_TMP=${OWN_TMP:-0}
FORCE=${FORCE:-0}
DOCKER=${DOCKER:-docker}

here=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
repo_root=$(cd -- "$here/../.." && pwd)

say() { echo "clean-space: $*"; }

# run <description> <command...>: run it, or under DRY_RUN print it.
run() {
  local what=$1; shift
  if [[ "$DRY_RUN" == 1 ]]; then
    say "would $what:"
    printf '    '; printf '%q ' "$@"; echo
  else
    say "$what"
    "$@"
  fi
}

# list <description> <command...>: print a command's output indented, or (none).
list() {
  local what=$1; shift
  local out
  out=$("$@" 2>/dev/null || true)
  say "$what:"
  if [[ -n "$out" ]]; then sed 's/^/    /' <<<"$out"; else echo "    (none)"; fi
}

free_kb() { df -Pk "$1" | awk 'NR==2 {print $4}'; }
human_gb() { awk -v kb="$1" 'BEGIN {printf "%.1f GB", kb / 1024 / 1024}'; }

mkdir -p "$BUILD_DIR"
before=$(free_kb "$BUILD_DIR")
say "$(human_gb "$before") free on the filesystem holding $BUILD_DIR"
[[ "$DRY_RUN" == 1 ]] && say "DRY_RUN=1: nothing will be removed"

# ---- 1. local: this repository's build output ------------------------------

say "local build output under $BUILD_DIR"
shopt -s nullglob dotglob
tmp_entries=("$BUILD_DIR"/tmp/*)
if ((${#tmp_entries[@]})); then
  run "remove ${#tmp_entries[@]} entries under $BUILD_DIR/tmp" rm -rf -- "${tmp_entries[@]}"
fi

if [[ -z "${KEEP_RELEASES// /}" ]]; then
  say "no releases selected (KEEP_RELEASES empty): keeping every release build directory"
else
  for dir in "$BUILD_DIR"/*/; do
    name=$(basename "$dir")
    [[ "$name" == tmp ]] && continue
    keep=0
    for k in $KEEP_RELEASES; do [[ "$name" == "$k" ]] && keep=1; done
    if ((keep)); then
      say "keeping $dir (selected release)"
    else
      run "remove stale build directory $dir ($(du -sh "$dir" 2>/dev/null | cut -f1))" rm -rf -- "$dir"
    fi
  done
fi
shopt -u dotglob

if [[ "$CLEAN_EXAMPLE_PROVIDERS" == 1 ]]; then
  example_dirs=("$repo_root"/examples/providers/terraform-provider-junos-*/ "$repo_root"/terraform-provider-junos-*/)
  for dir in "${example_dirs[@]}"; do
    run "remove example provider output $dir" rm -rf -- "$dir"
  done
fi

# ---- 2. containers: only what Docker itself calls prunable -----------------

docker_ok=0
if command -v "$DOCKER" >/dev/null 2>&1 && "$DOCKER" info >/dev/null 2>&1; then
  docker_ok=1
fi

if ((docker_ok)); then
  say "containers (runtime: $DOCKER)"
  running=$("$DOCKER" ps -q | wc -l)
  say "$running running container(s) will not be touched"
  if [[ "$DRY_RUN" == 1 ]]; then
    list "exited containers older than $PRUNE_UNTIL that would be removed" \
      "$DOCKER" ps -a --filter status=exited --filter status=dead --filter status=created \
      --format '{{.ID}}  {{.Names}}  {{.Status}}  ({{.Image}})'
    list "dangling images that would be removed" \
      "$DOCKER" images -f dangling=true --format '{{.ID}}  {{.Repository}}:{{.Tag}}  {{.Size}}'
    if [[ "$PRUNE_UNUSED_IMAGES" == 1 ]]; then
      list "tagged images (PRUNE_UNUSED_IMAGES=1 removes those no container uses${PRUNE_IMAGE_AGE:+, older than $PRUNE_IMAGE_AGE})" \
        "$DOCKER" images --format '{{.ID}}  {{.Repository}}:{{.Tag}}  {{.Size}}  {{.CreatedSince}}'
    fi
    list "docker system df" "$DOCKER" system df
  fi
  run "prune stopped containers older than $PRUNE_UNTIL" "$DOCKER" container prune -f --filter "until=$PRUNE_UNTIL"
  run "prune dangling images" "$DOCKER" image prune -f
  run "prune unreferenced build cache older than $PRUNE_UNTIL" "$DOCKER" builder prune -f --filter "until=$PRUNE_UNTIL"
  if [[ "$PRUNE_UNUSED_IMAGES" == 1 ]]; then
    if [[ -n "$PRUNE_IMAGE_AGE" ]]; then
      run "prune tagged images unused by any container and older than $PRUNE_IMAGE_AGE" \
        "$DOCKER" image prune -a -f --filter "until=$PRUNE_IMAGE_AGE"
    else
      run "prune every tagged image unused by any container" "$DOCKER" image prune -a -f
    fi
  fi
else
  say "container cleanup skipped: $DOCKER is not available"
fi

# ---- 3. opt-in: Go caches, own /tmp ----------------------------------------

if [[ "$CLEAN_GO_CACHE" == 1 ]]; then
  if command -v go >/dev/null 2>&1; then
    run "remove the Go build cache ($(du -sh "$(go env GOCACHE)" 2>/dev/null | cut -f1))" go clean -cache
  else
    say "CLEAN_GO_CACHE=1 but go is not available"
  fi
fi

if [[ "$CLEAN_GO_MODCACHE" == 1 ]]; then
  if command -v go >/dev/null 2>&1; then
    modcache=$(go env GOMODCACHE)
    if [[ "$DRY_RUN" != 1 && "$FORCE" != 1 ]]; then
      read -r -p "clean-space: remove the Go module cache $modcache ($(du -sh "$modcache" 2>/dev/null | cut -f1))? Every Go project re-downloads. [y/N] " answer
      [[ "$answer" == [yY]* ]] || { say "keeping the Go module cache"; CLEAN_GO_MODCACHE=0; }
    fi
    if [[ "$CLEAN_GO_MODCACHE" == 1 ]]; then
      run "remove the Go module cache $modcache" go clean -modcache
    fi
  else
    say "CLEAN_GO_MODCACHE=1 but go is not available"
  fi
fi

if [[ "$OWN_TMP" == 1 ]]; then
  tmp_root=${TMP_ROOT:-/tmp}
  mapfile -t own < <(find "$tmp_root" -mindepth 1 -maxdepth 1 -user "$(id -u)" -mtime +7 2>/dev/null)
  if ((${#own[@]})); then
    run "remove ${#own[@]} of this user's entries under $tmp_root older than 7 days" rm -rf -- "${own[@]}"
  else
    say "nothing of this user's under $tmp_root is older than 7 days"
  fi
fi

# ---- report ----------------------------------------------------------------

after=$(free_kb "$BUILD_DIR")
if [[ "$DRY_RUN" == 1 ]]; then
  say "dry run complete; nothing removed"
else
  say "reclaimed $(human_gb $((after - before))) on the filesystem holding $BUILD_DIR ($(human_gb "$after") free now)"
fi
