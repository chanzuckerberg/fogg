# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

TF_VARS := $(patsubst %,-e%,$(filter TF_VAR_%,$(.VARIABLES)))
REPO_ROOT := $(shell git rev-parse --show-toplevel)
REPO_RELATIVE_PATH := $(shell git rev-parse --show-prefix)
AUTO_APPROVE := false
# We need to do this because `terraform fmt` recurses into .terraform/modules
# and wont' accept more than one file at a time.
TF=$(wildcard *.tf)
IMAGE_VERSION=$(DOCKER_IMAGE_VERSION)_TF$(TERRAFORM_VERSION)

# dependencies.mk helps with dependencies such as terraform
include: ./dependencies.mk

docker_base := \
	docker run --rm -e HOME=/home -v $$HOME/.aws:/home/.aws -v $(REPO_ROOT):/repo \
	-v $(REPO_ROOT)/.bin:/usr/local/bin -v $(REPO_ROOT)/terraform.d:/repo/$(REPO_RELATIVE_PATH)/terraform.d \
	-e GIT_SSH_COMMAND='ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no' \
	-e RUN_USER_ID=$(shell id -u) -e RUN_GROUP_ID=$(shell id -g) \
	-e TF_PLUGIN_CACHE_DIR="/repo/.terraform.d/plugin-cache" -e TF="$(TF)" \
	-w /repo/$(REPO_RELATIVE_PATH) $(TF_VARS) $(FOGG_DOCKER_FLAGS) $$(sh $(REPO_ROOT)/scripts/docker-ssh-mount.sh)
docker_terraform := $(docker_base) chanzuckerberg/terraform:$(IMAGE_VERSION)
docker_sh := $(docker_base) --entrypoint='/bin/sh' chanzuckerberg/terraform:$(IMAGE_VERSION)

ifdef FOGG_DISABLE_DOCKER
	sh_command ?= $(SHELL)
	terraform_command ?= terraform
else
	sh_command ?= $(docker_sh)
	terraform_command ?= $(docker_terraform)
endif

all:

fmt:
	@$(sh_command) -c 'for f in $(TF); do printf .; terraform fmt $$f; done'; \
	echo

lint: terraform-validate lint-terraform-fmt lint-tflint

lint-tflint: init
	@if (( $$TFLINT_ENABLED )); then \
    $(sh_command) -c 'tflint' || exit $$?; \
	else \
    echo "tflint not enabled"; \
	fi \

terraform-validate: init
	@$(sh_command) -c 'terraform validate -check-variables=true $$f || exit $$?'

lint-terraform-fmt:
	@$(sh_command) -c 'for f in $(TF); do printf .; terraform fmt --check=true --diff=true $$f || exit $$? ; done'

get: ssh-forward
	$(terraform_command) get --update=true

plan: fmt get init ssh-forward
	$(terraform_command) plan

apply: FOGG_DOCKER_FLAGS = -it
apply: fmt get init ssh-forward
	$(terraform_command) apply -auto-approve=$(AUTO_APPROVE)

docs:
	@echo

clean:
	-rm -rfv .terraform/modules
	-rm -rfv .terraform/plugins

test:

init: ssh-forward
	$(terraform_command) init -input=false

check-plan: init get ssh-forward
	$(terraform_command) plan -detailed-exitcode; \
	ERR=$$?; \
	if [ $$ERR -eq 0 ] ; then \
		echo "Success"; \
	elif [ $$ERR -eq 1 ] ; then \
		echo "Error in plan execution."; \
		exit 1; \
	elif [ $$ERR -eq 2 ] ; then \
		echo "Diff";  \
	fi

ssh-forward:
ifndef FOGG_DISABLE_DOCKER
	bash $(REPO_ROOT)/scripts/docker-ssh-forward.sh
endif

run: FOGG_DOCKER_FLAGS = -it
run:
	$(terraform_command) $(CMD)

.PHONY: all apply clean docs fmt get lint plan run ssh-forward test
