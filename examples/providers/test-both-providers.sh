#!/usr/bin/env bash

set -euo pipefail

providers_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "$providers_dir/../.." && pwd)"
terraform_source_dir="${JTAF_TERRAFORM_SOURCE_DIR:-$repo_root/examples/terraform_files}"
stage_root="$(mktemp -d "${TMPDIR:-/tmp}/jtaf-provider-e2e.XXXXXX")"
terraform_dir="$stage_root/terraform"
standard_dir="$stage_root/standard"
generic_dir="$stage_root/generic"

cleanup() {
  rm -rf "$stage_root"
}
trap cleanup EXIT

for command in go terraform jtaf-yang2go jtaf-xml2tf; do
  if ! command -v "$command" >/dev/null 2>&1; then
    echo "Required command not found: $command" >&2
    exit 1
  fi
done

if [[ ! -d "$terraform_source_dir" ]]; then
  echo "Terraform source directory not found: $terraform_source_dir" >&2
  exit 1
fi

mkdir -p "$standard_dir" "$generic_dir" "$terraform_dir"
cp -a "$terraform_source_dir"/. "$terraform_dir"/

echo "Generating standard provider variants..."
JTAF_PROVIDER_OUTPUT_DIR="$standard_dir" bash "$providers_dir/build.sh"

echo "Generating generic provider variants..."
JTAF_PROVIDER_OUTPUT_DIR="$generic_dir" \
JTAF_SKIP_PROVIDER_INSTALL=1 \
  bash "$providers_dir/build-generic.sh"

for provider_type in vqfx-evpn-vxlan vsrx-evpn-vxlan; do
  standard_module="$standard_dir/terraform-provider-junos-$provider_type"
  generic_module="$generic_dir/terraform-provider-junos-$provider_type"
  for module in "$standard_module" "$generic_module"; do
    if [[ ! -s "$module/trimmed_schema.json.gz" ]]; then
      echo "Expected schema not found: $module/trimmed_schema.json.gz" >&2
      exit 1
    fi
  done

  if [[ ! -x "$standard_module/terraform-provider-junos-$provider_type" ]]; then
    (cd "$standard_module" && go build .)
  fi
  for module in "$standard_module" "$generic_module"; do
    if [[ ! -x "$module/terraform-provider-junos-$provider_type" ]]; then
      echo "Expected provider binary not found: $module/terraform-provider-junos-$provider_type" >&2
      exit 1
    fi
  done
done

echo "Generating the shared Terraform scenario..."
JTAF_PROVIDER_OUTPUT_DIR="$standard_dir" \
JTAF_TERRAFORM_OUTPUT_DIR="$terraform_dir" \
  bash "$providers_dir/convert.sh"
rm -f \
  "$terraform_dir/preview.plan" \
  "$terraform_dir/preview_generic.plan" \
  "$terraform_dir/preview_full_config.txt" \
  "$terraform_dir/preview_generic_full_config.txt"

hcl_string() {
  local escaped="${1//\\/\\\\}"
  escaped="${escaped//\"/\\\"}"
  printf '"%s"' "$escaped"
}

write_cli_config() {
  local config_path="$1"
  local variant_dir="$2"
  {
    printf 'provider_installation {\n  dev_overrides {\n'
    printf '    "registry.terraform.io/hashicorp/junos-vqfx-evpn-vxlan" = %s\n' \
      "$(hcl_string "$variant_dir/terraform-provider-junos-vqfx-evpn-vxlan")"
    printf '    "registry.terraform.io/hashicorp/junos-vsrx-evpn-vxlan" = %s\n' \
      "$(hcl_string "$variant_dir/terraform-provider-junos-vsrx-evpn-vxlan")"
    printf '  }\n  direct {}\n}\n'
  } > "$config_path"
}

write_cli_config "$stage_root/terraformrc-standard" "$standard_dir"
write_cli_config "$stage_root/terraformrc-generic" "$generic_dir"

git_commit="$(git -C "$repo_root" rev-parse --short HEAD 2>/dev/null || printf unknown)"
generated_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

run_plan() {
  local variant="$1"
  local config_path="$2"
  local plan_path="$3"
  local output_path="$4"
  local raw_output="$stage_root/$variant-plan.log"
  local exit_code=0

  echo "Running Terraform plan with the $variant provider..."
  if (cd "$terraform_dir" && TF_CLI_CONFIG_FILE="$config_path" \
      terraform plan -no-color -out="$plan_path") >"$raw_output" 2>&1; then
    exit_code=0
  else
    exit_code=$?
  fi

  if [[ "$exit_code" -eq 0 && ! -s "$plan_path" ]]; then
    printf 'Terraform plan succeeded without creating %s.\n' "$plan_path" >> "$raw_output"
    exit_code=1
  fi

  if [[ "$exit_code" -eq 0 ]]; then
    {
      printf '\n# planned_configuration\n'
      (cd "$terraform_dir" && TF_CLI_CONFIG_FILE="$config_path" \
        terraform show -no-color "$plan_path")
    } >> "$raw_output" 2>&1 || exit_code=$?
  fi

  {
    printf '# preview_metadata\n'
    printf 'generated_at=%s\n' "$generated_at"
    printf 'provider_variant=%s\n' "$variant"
    printf 'git_commit=%s\n' "$git_commit"
    printf 'plan_exit_code=%s\n\n' "$exit_code"
    cat "$raw_output"
  } > "$output_path"
  cat "$raw_output"
  printf '%s_plan_exit_code=%s\n' "$variant" "$exit_code"
  return "$exit_code"
}

standard_exit_code=0
generic_exit_code=0
run_plan standard "$stage_root/terraformrc-standard" \
  "$terraform_dir/preview.plan" "$terraform_dir/preview_full_config.txt" \
  || standard_exit_code=$?
run_plan generic "$stage_root/terraformrc-generic" \
  "$terraform_dir/preview_generic.plan" "$terraform_dir/preview_generic_full_config.txt" \
  || generic_exit_code=$?

for artifact in preview_full_config.txt preview_generic_full_config.txt; do
  if [[ ! -s "$terraform_dir/$artifact" ]]; then
    echo "Expected plan result is missing or empty: $terraform_dir/$artifact" >&2
    exit 1
  fi
  cp "$terraform_dir/$artifact" "$terraform_source_dir/$artifact"
done

if [[ "$standard_exit_code" -eq 0 && -s "$terraform_dir/preview.plan" ]]; then
  cp "$terraform_dir/preview.plan" "$terraform_source_dir/preview.plan"
else
  rm -f "$terraform_source_dir/preview.plan"
fi
if [[ "$generic_exit_code" -eq 0 && -s "$terraform_dir/preview_generic.plan" ]]; then
  cp "$terraform_dir/preview_generic.plan" "$terraform_source_dir/preview_generic.plan"
else
  rm -f "$terraform_source_dir/preview_generic.plan"
fi

printf '\nstandard_plan_exit_code=%s\ngeneric_plan_exit_code=%s\n' \
  "$standard_exit_code" "$generic_exit_code"
if [[ "$standard_exit_code" -ne 0 || "$generic_exit_code" -ne 0 ]]; then
  exit 1
fi