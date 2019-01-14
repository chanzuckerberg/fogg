# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

TF_VARS := $(patsubst %,-e%,$(filter TF_VAR_%,$(.VARIABLES)))
REPO_ROOT := $(shell git rev-parse --show-toplevel)
REPO_RELATIVE_PATH := $(shell git rev-parse --show-prefix)
AUTO_APPROVE := false
# We need to do this because `terraform fmt` recurses into .terraform/modules
# and wont' accept more than one file at a time.
TF=$(wildcard *.tf)


ifdef USE_DOCKER
	IMAGE_VERSION=$(DOCKER_IMAGE_VERSION)_TF$(TERRAFORM_VERSION)
	docker_base = \
		docker run --rm -e HOME=/home -v $$HOME/.aws:/home/.aws -v $(REPO_ROOT):/repo \
		-v $(REPO_ROOT)/.fogg/bin:/usr/local/bin -v $(REPO_ROOT)/terraform.d:/repo/$(REPO_RELATIVE_PATH)/terraform.d \
		-e GIT_SSH_COMMAND='ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no' \
		-e RUN_USER_ID=$(shell id -u) -e RUN_GROUP_ID=$(shell id -g) \
		-e TF_PLUGIN_CACHE_DIR="/repo/.terraform.d/plugin-cache" -e TF="$(TF)" \
		-w /repo/$(REPO_RELATIVE_PATH) $(TF_VARS) $(FOGG_DOCKER_FLAGS) $$(sh $(REPO_ROOT)/scripts/docker-ssh-mount.sh)
	docker_terraform = $(docker_base) chanzuckerberg/terraform:$(IMAGE_VERSION)
	docker_sh = $(docker_base) --entrypoint='/bin/sh' chanzuckerberg/terraform:$(IMAGE_VERSION)
	sh_command ?= $(docker_sh)
	terraform_command ?= $(docker_terraform)
else
	TFENV_DIR ?= $(HOME)/.tfenv
	export PATH :=$(TFENV_DIR)/versions/$(TERRAFORM_VERSION)/:$(PATH)
	export TF_PLUGIN_CACHE_DIR=$(REPO_ROOT)/.terraform.d/plugin-cache
	sh_command ?= $(SHELL)
	terraform_command ?= $(TFENV_DIR)/versions/$(TERRAFORM_VERSION)/terraform
endif

tfenv:
ifndef USE_DOCKER
	if [ ! -d ${TFENV_DIR} ]; then \
		echo; \
		git clone https://github.com/tfutils/tfenv.git $(TFENV_DIR); \
	fi
endif
.PHONY: tfenv

terraform: tfenv
ifndef USE_DOCKER
	${TFENV_DIR}/bin/tfenv install $(TERRAFORM_VERSION)
endif
.PHONY: terraform
