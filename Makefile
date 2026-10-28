# Build the generic, groups-aware Junos Terraform providers for a Junos OS and
# a Junos OS Evolved release, one provider per family, from the Juniper YANG
# repository (https://github.com/Juniper/yang).
#
#   make latest            which releases would be built, from which commit
#   make all               every provider for both releases -> dist/<release>/
#   make junos | evo       one branch
#   make qfx               one family (mx es ex nfx qfx, or evo)
#   make clean-space       reclaim unused space on the build machine (DRY_RUN=1 first)
#   make help              targets and variables with their current values
#
# The build does not touch the YANG clone's working tree: modules are read
# out of its git objects at YANG_REF. Outputs for different releases sit side
# by side under dist/, each with a manifest and a Terraform CLI configuration
# whose dev_overrides point at it.

SHELL := /bin/bash
.SHELLFLAGS := -e -o pipefail -c
.DELETE_ON_ERROR:
.DEFAULT_GOAL := help
MAKEFLAGS += --no-builtin-rules --no-print-directory

# ---------------------------------------------------------------------------
# Variables (override on the command line: make all JOBS=1 EXCLUDE="tenants")
# ---------------------------------------------------------------------------

YANG_REPO ?= $(abspath ../yang)
YANG_REF ?= HEAD
YANG_DIR ?=
BUILD_DIR ?= build/release
DIST_DIR ?= dist
JOBS ?= 2
NICENESS ?= 10
MIN_FREE_GB ?= 5
EXCLUDE ?=
# Per-family exclusions, added to EXCLUDE for that provider only: EXCLUDE_mx,
# EXCLUDE_es, EXCLUDE_ex, EXCLUDE_nfx, EXCLUDE_qfx, EXCLUDE_evo.
#
# NFX: the hash-function choice under enhanced-hash-key has both a leaf CRC16
# and a container crc16, which map to the same Terraform attribute, so
# jtaf-provider refuses the model (25.4R1; the other families have the leaf
# alone). Until the generator can tell them apart, the choice is left out of
# the NFX provider at each place it occurs, in the base hierarchy and the
# groups copy. Set EXCLUDE_nfx= to clear it for a release that does not need it.
NFX_HASH_FUNCTION := forwarding-options/enhanced-hash-key/hash-function
NFX_HASH_FUNCTION_PATHS := $(NFX_HASH_FUNCTION) \
  routing-instances/instance/$(NFX_HASH_FUNCTION) \
  logical-systems/routing-instances/instance/$(NFX_HASH_FUNCTION) \
  tenants/routing-instances/instance/$(NFX_HASH_FUNCTION)
EXCLUDE_nfx ?= $(NFX_HASH_FUNCTION_PATHS) $(addprefix groups/,$(NFX_HASH_FUNCTION_PATHS))
YANG2GO_FLAGS ?=
FAMILIES ?=
PROVIDER_NAMESPACE ?= hashicorp
AUTO_CLEAN ?= 0
# JUNOS_RELEASE / EVO_RELEASE default to the newest release of each branch at
# YANG_REF; set one to empty (EVO_RELEASE=) to leave that branch out.

ifeq ($(strip $(JOBS)),)
  $(error JOBS must be a number)
endif
MAKEFLAGS += -j$(JOBS)

RELEASE_PY := python3 scripts/release/jtaf_release.py
ifneq ($(YANG_DIR),)
  TREE_ARGS := --dir $(YANG_DIR)
  YANG_SOURCE := dir:$(abspath $(YANG_DIR))
else
  TREE_ARGS := --repo $(YANG_REPO) --ref $(YANG_REF)
  YANG_SOURCE := git:$(YANG_REPO)@$(YANG_REF)
endif

# Everything the build and the Go toolchain write temporarily goes under the
# build directory, not the system temporary directory.
export TMPDIR := $(abspath $(BUILD_DIR))/tmp
export GOTMPDIR := $(TMPDIR)

NICE := nice -n $(NICENESS)$(if $(shell command -v ionice 2>/dev/null), ionice -c 2 -n 7,)
exclude_flags = $(foreach e,$(EXCLUDE) $(EXCLUDE_$(1)),--exclude $(e))

UTILITY_GOALS := help latest clean distclean clean-space preflight
BUILD_GOALS := $(filter-out $(UTILITY_GOALS),$(MAKECMDGOALS))

# ---------------------------------------------------------------------------
# Release resolution, only when a build goal asks for it
# ---------------------------------------------------------------------------

ifneq ($(BUILD_GOALS),)

ifeq ($(origin JUNOS_RELEASE),undefined)
  JUNOS_RELEASE := $(shell $(RELEASE_PY) $(TREE_ARGS) -q latest --branch junos)
endif
ifeq ($(origin EVO_RELEASE),undefined)
  EVO_RELEASE := $(shell $(RELEASE_PY) $(TREE_ARGS) -q latest --branch evo)
endif
ifeq ($(strip $(JUNOS_RELEASE)$(EVO_RELEASE)),)
  $(error no release to build: could not discover one in $(YANG_SOURCE), and none was given)
endif
ifeq ($(origin YANG_COMMIT),undefined)
  YANG_COMMIT := $(shell $(RELEASE_PY) $(TREE_ARGS) -q commit)
endif
ifeq ($(strip $(YANG_COMMIT)),)
  $(error cannot resolve $(YANG_SOURCE))
endif
ifeq ($(origin TOOLING_ID),undefined)
  TOOLING_ID := $(shell $(RELEASE_PY) tooling-id)
endif

# Each release's layout (module paths, families) is written once per commit
# as an included makefile; make generates it and restarts.
RELEASES := $(strip $(JUNOS_RELEASE) $(EVO_RELEASE))
-include $(foreach r,$(RELEASES),$(BUILD_DIR)/$(r)/release-$(YANG_COMMIT).mk)

$(BUILD_DIR)/%/release-$(YANG_COMMIT).mk:
	@mkdir -p $(@D)
	@$(RELEASE_PY) $(TREE_ARGS) -q makevars $* --prefix $(if $(filter %-EVO,$*),EVO,JUNOS) > $@.tmp
	@mv $@.tmp $@

ifneq ($(MAKE_RESTARTS),)
  ifneq ($(JUNOS_RELEASE),)
    ifeq ($(JUNOS_TYPES),)
      $(error could not read the layout of $(JUNOS_RELEASE) from $(YANG_SOURCE))
    endif
  endif
  ifneq ($(EVO_RELEASE),)
    ifeq ($(EVO_TYPES),)
      $(error could not read the layout of $(EVO_RELEASE) from $(YANG_SOURCE))
    endif
  endif
endif

JUNOS_FAMILIES := $(or $(FAMILIES),$(JUNOS_TYPES))
ifneq ($(filter-out $(JUNOS_TYPES),$(JUNOS_FAMILIES)),)
  ifneq ($(JUNOS_TYPES),)
    $(error FAMILIES has $(filter-out $(JUNOS_TYPES),$(JUNOS_FAMILIES)); $(JUNOS_RELEASE) has $(JUNOS_TYPES))
  endif
endif

endif # BUILD_GOALS

# ---------------------------------------------------------------------------
# Rules
# ---------------------------------------------------------------------------

# Each stage depends on a stamp whose name carries a hash of the stage's
# inputs (computed when make reads this file), so changed inputs mean a new,
# newer stamp and the stage reruns; unchanged inputs mean nothing to do. No
# recipe has to run to find that out, so `make -n` reports accurately.
SHA := $(if $(shell command -v sha256sum 2>/dev/null),sha256sum,shasum -a 256)
inputs_key = $(shell printf '%s\n' $(1) | $(SHA) | cut -c1-12)

# ---- export: a release's common and configuration modules -----------------

# $(1) release
define export_rules
EXPORT_INPUTS_$(1) := 'release=$(1)' 'source=$(YANG_SOURCE)' 'commit=$(YANG_COMMIT)'
EXPORT_KEY_$(1) := $$(call inputs_key,$$(EXPORT_INPUTS_$(1)))

$(BUILD_DIR)/$(1)/yang/.inputs-$$(EXPORT_KEY_$(1)):
	@mkdir -p $$(@D)
	@printf '%s\n' $$(EXPORT_INPUTS_$(1)) > $$@

$(BUILD_DIR)/$(1)/yang/.exported: $(BUILD_DIR)/$(1)/yang/.inputs-$$(EXPORT_KEY_$(1))
	@echo "[$(1)] exporting modules from $(YANG_SOURCE) ($(YANG_COMMIT))"
	@find $$(@D) -mindepth 1 -maxdepth 1 ! -name '$$(<F)' -exec rm -rf {} +
	@$(RELEASE_PY) $(TREE_ARGS) -q export $(1) $$(@D)
	@cp $$< $$(@D)/inputs.txt
	@touch $$@
endef

# ---- generate, compile, distribute: one provider per family ---------------

# $(1) release  $(2) variable prefix (JUNOS|EVO)  $(3) provider type
define provider_rules
GEN_INPUTS_$(1)_$(3) := 'release=$(1)' 'commit=$(YANG_COMMIT)' 'type=$(3)' 'modules=$($(2)_COMMON) $($(2)_CONF_$(3))' 'exclude=$(strip $(EXCLUDE) $(EXCLUDE_$(3)))' 'flags=$(YANG2GO_FLAGS)' 'tooling=$(TOOLING_ID)'
GEN_KEY_$(1)_$(3) := $$(call inputs_key,$$(GEN_INPUTS_$(1)_$(3)))

$(BUILD_DIR)/$(1)/.inputs-$(3)-$$(GEN_KEY_$(1)_$(3)):
	@mkdir -p $$(@D)
	@printf '%s\n' $$(GEN_INPUTS_$(1)_$(3)) > $$@

$(BUILD_DIR)/$(1)/.generated-$(3): $(BUILD_DIR)/$(1)/yang/.exported $(BUILD_DIR)/$(1)/.inputs-$(3)-$$(GEN_KEY_$(1)_$(3)) | preflight
	@echo "[$(1)] generating junos-$(3): --generic --groups $(call exclude_flags,$(3)) $(YANG2GO_FLAGS)"
	@mkdir -p $(TMPDIR)
	@cd $(BUILD_DIR)/$(1) && $(NICE) jtaf-yang2go --generic --groups $(call exclude_flags,$(3)) $(YANG2GO_FLAGS) \
	  -p yang/$($(2)_COMMON) yang/$($(2)_CONF_$(3))/*.yang -t $(3) \
	  > generate-$(3).log 2>&1 || { cat generate-$(3).log; exit 1; }
	@cp $$(word 2,$$^) $(BUILD_DIR)/$(1)/inputs-$(3).txt
	@find $$(@D) -maxdepth 1 -name '.inputs-$(3)-*' ! -name '$$(notdir $$(word 2,$$^))' -delete
	@touch $$@

$(BUILD_DIR)/$(1)/terraform-provider-junos-$(3)/terraform-provider-junos-$(3): $(BUILD_DIR)/$(1)/.generated-$(3)
	@echo "[$(1)] compiling junos-$(3)"
	@mkdir -p $(TMPDIR)
	@cd $$(@D) && $(NICE) go build .

$(DIST_DIR)/$(1)/bin/terraform-provider-junos-$(3): $(BUILD_DIR)/$(1)/terraform-provider-junos-$(3)/terraform-provider-junos-$(3)
	@mkdir -p $$(@D) $(DIST_DIR)/$(1)/schemas
	@cp $$< $$@
	@cp $$(<D)/trimmed_schema.json.gz $(DIST_DIR)/$(1)/schemas/junos-$(3).trimmed_schema.json.gz
	@echo "[$(1)] junos-$(3) -> $$@"

.PHONY: $(3)
$(3): $(DIST_DIR)/$(1)/bin/terraform-provider-junos-$(3)
	$$(call release_summary,$(1))
endef

# Manifest and Terraform CLI configuration for whatever a release has built.
define release_summary
@$(RELEASE_PY) -q manifest $(BUILD_DIR)/$(1) $(DIST_DIR)/$(1) --source '$(YANG_SOURCE)' --commit $(YANG_COMMIT) --tooling-id $(TOOLING_ID)
@$(RELEASE_PY) -q terraformrc $(DIST_DIR)/$(1) --namespace $(PROVIDER_NAMESPACE)
endef

$(foreach r,$(RELEASES),$(eval $(call export_rules,$(r))))

ifneq ($(JUNOS_RELEASE),)
$(foreach t,$(JUNOS_TYPES),$(eval $(call provider_rules,$(JUNOS_RELEASE),JUNOS,$(t))))
JUNOS_BINS := $(foreach t,$(JUNOS_FAMILIES),$(DIST_DIR)/$(JUNOS_RELEASE)/bin/terraform-provider-junos-$(t))

.PHONY: junos
junos: $(JUNOS_BINS)
	$(call release_summary,$(JUNOS_RELEASE))
	@echo "[$(JUNOS_RELEASE)] done: $(JUNOS_FAMILIES)"
BRANCHES += junos
endif

ifneq ($(EVO_RELEASE),)
$(foreach t,$(EVO_TYPES),$(eval $(call provider_rules,$(EVO_RELEASE),EVO,$(t))))
EVO_BINS := $(foreach t,$(EVO_TYPES),$(DIST_DIR)/$(EVO_RELEASE)/bin/terraform-provider-junos-$(t))
ifneq ($(filter evo,$(EVO_TYPES)),evo)
.PHONY: evo
evo: $(EVO_BINS)
	$(call release_summary,$(EVO_RELEASE))
endif
BRANCHES += evo
endif

.PHONY: all
all: $(BRANCHES)
	@$(RELEASE_PY) -q terraformrc --combined $(DIST_DIR) --namespace $(PROVIDER_NAMESPACE)
	@echo "done: $(RELEASES) -> $(DIST_DIR)/  (TF_CLI_CONFIG_FILE=$(abspath $(DIST_DIR))/terraformrc)"

# ---------------------------------------------------------------------------
# Utilities
# ---------------------------------------------------------------------------

.PHONY: latest
latest:
	@$(RELEASE_PY) $(TREE_ARGS) -q latest

.PHONY: preflight
preflight:
	@BUILD_DIR='$(BUILD_DIR)' MIN_FREE_GB='$(MIN_FREE_GB)' JOBS='$(JOBS)' AUTO_CLEAN='$(AUTO_CLEAN)' \
	  bash scripts/release/preflight.sh

.PHONY: clean-space
# The releases currently selected keep their build directories; if none were
# named, the newest of each branch is looked up, best effort (no clone, no
# release build directory is removed).
clean-space:
	@keep="$(JUNOS_RELEASE) $(EVO_RELEASE)"; \
	if [ -z "$${keep// /}" ]; then \
	  keep="$$($(RELEASE_PY) $(TREE_ARGS) -q latest --branch junos 2>/dev/null || true) $$($(RELEASE_PY) $(TREE_ARGS) -q latest --branch evo 2>/dev/null || true)"; \
	fi; \
	BUILD_DIR='$(BUILD_DIR)' KEEP_RELEASES="$$keep" bash scripts/release/clean_space.sh

.PHONY: clean distclean
clean:
	$(if $(strip $(BUILD_DIR)),rm -rf '$(BUILD_DIR)',@echo "BUILD_DIR is empty")

distclean: clean
	$(if $(strip $(DIST_DIR)),rm -rf '$(DIST_DIR)',@echo "DIST_DIR is empty")

.PHONY: help
help:
	@echo "Junos provider release matrix"
	@echo
	@echo "Targets:"
	@echo "  latest        releases that would be built, and the YANG commit"
	@echo "  all           every provider for both releases -> \$$(DIST_DIR)/<release>/"
	@echo "  junos | evo   one branch"
	@echo "  mx es ex nfx qfx | evo   one family"
	@echo "  preflight     check tools, disk and memory before a build"
	@echo "  clean-space   reclaim unused space on the build machine"
	@echo "                  DRY_RUN=1 preview; PRUNE_UNTIL=24h; opt-in: PRUNE_UNUSED_IMAGES=1 [PRUNE_IMAGE_AGE=168h],"
	@echo "                  CLEAN_GO_CACHE=1, CLEAN_GO_MODCACHE=1 [FORCE=1], OWN_TMP=1, CLEAN_EXAMPLE_PROVIDERS=1"
	@echo "  clean         remove \$$(BUILD_DIR);  distclean: also \$$(DIST_DIR)"
	@echo "  help          this text (the default goal)"
	@echo
	@echo "Variables (current values):"
	@echo "  YANG_REPO=$(YANG_REPO)  YANG_REF=$(YANG_REF)  YANG_DIR=$(YANG_DIR)"
	@echo "  JUNOS_RELEASE / EVO_RELEASE   default: newest of each branch (empty skips the branch)"
	@echo "  FAMILIES=$(FAMILIES)   default: every family the release has"
	@echo "  EXCLUDE=$(EXCLUDE)   subtrees left out of every provider, e.g. \"logical-systems tenants dynamic-profiles\""
	@echo "  EXCLUDE_<type>   subtrees left out of one provider; EXCLUDE_nfx defaults to the hash-function"
	@echo "                   choice (CRC16/crc16 clash in the 25.4R1 NFX model): $(words $(EXCLUDE_nfx)) paths"
	@echo "  YANG2GO_FLAGS=$(YANG2GO_FLAGS)   extra jtaf-yang2go flags"
	@echo "  JOBS=$(JOBS)  NICENESS=$(NICENESS)  MIN_FREE_GB=$(MIN_FREE_GB)  AUTO_CLEAN=$(AUTO_CLEAN)"
	@echo "  BUILD_DIR=$(BUILD_DIR)  DIST_DIR=$(DIST_DIR)  PROVIDER_NAMESPACE=$(PROVIDER_NAMESPACE)"
	@echo
	@echo "Each release's providers embed the whole model with configuration groups (--generic --groups);"
	@echo "planning with one needs several GB. Point Terraform at a set with"
	@echo "  TF_CLI_CONFIG_FILE=\$$PWD/$(DIST_DIR)/terraformrc"
