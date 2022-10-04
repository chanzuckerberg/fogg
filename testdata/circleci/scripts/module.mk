# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

SELF_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(SELF_DIR)/common.mk

all: fmt lint doc
.PHONY: all

fmt: terraform ## run terraform fmt on this module
	@$(terraform_command) fmt $(TF_ARGS)
.PHONY: fmt

validate: terraform ## run terraform fmt on this module
	@$(terraform_command) validate $(TF_ARGS)
.PHONY: validate

check: lint check-docs ## run all checks on this module
.PHONY: check

lint: lint-tf check-docs ## run all linters on this module
.PHONY: lint

lint-tf: terraform ## run terraform linters on this module
	$(terraform_command) fmt $(TF_ARGS) --check=true --diff=true
.PHONY: lint-tf

readme: ## update this module's README.md
	bash $(REPO_ROOT)/scripts/update-readme.sh update
.PHONY: readme

docs: readme ## update all docs for this module
.PHONY: docs

check-docs: ## check that this module's docs are up-to-date
	@bash $(REPO_ROOT)/scripts/update-readme.sh check; \
	if [ ! $$? -eq 0 ];  then \
		echo "Docs are out of date, run \`make docs\`"; \
		exit 1 ; \
	fi
.PHONY: check-docs

clean:
.PHONY: clean

test:
.PHONY: test
