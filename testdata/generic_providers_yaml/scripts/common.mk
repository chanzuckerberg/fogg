# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

SHELL=/bin/bash -o pipefail
TF_VARS := $(patsubst %,-e%,$(filter TF_VAR_%,$(.VARIABLES)))
REPO_ROOT := $(shell git rev-parse --show-toplevel)
REPO_RELATIVE_PATH := $(shell git rev-parse --show-prefix)
AUTO_APPROVE := false
export AWS_SDK_LOAD_CONFIG := 1

TFENV_DIR ?= $(REPO_ROOT)/.fogg/tfenv
export PATH :=$(TFENV_DIR)/libexec:$(TFENV_DIR)/versions/$(TERRAFORM_VERSION)/:$(REPO_ROOT)/.fogg/bin:$(PATH)
export TF_PLUGIN_CACHE_DIR=$(REPO_ROOT)/.terraform.d/plugin-cache
export TF_IN_AUTOMATION=1
terraform_command ?= $(TFENV_DIR)/versions/$(TERRAFORM_VERSION)/terraform

ifeq ($(TF_BACKEND_KIND),remote)
	REFRESH := true
else
	REFRESH := false
endif


tfenv: ## install the tfenv tool
	@if [ ! -d ${TFENV_DIR} ]; then \
		git clone -q https://github.com/tfutils/tfenv.git $(TFENV_DIR); \
	fi
.PHONY: tfenv

terraform: tfenv ## ensure that the proper version of terraform is installed
	@${TFENV_DIR}/bin/tfenv install $(TERRAFORM_VERSION)
.PHONY: terraform