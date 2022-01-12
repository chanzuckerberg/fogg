SHA=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)
DIRTY=false
# TODO add release flag
GO_PACKAGE=$(shell go list)
LDFLAGS=-ldflags "-w -s -X $(GO_PACKAGE)/util.GitSha=${SHA} -X $(GO_PACKAGE)/util.Version=${VERSION} -X $(GO_PACKAGE)/util.Dirty=${DIRTY}"
export GO111MODULE=on

all: test install

fmt:
	goimports -w -l .
.PHONY: fmt

lint-setup: ## setup linter dependencies
	## See: https://github.com/igorshubovych/markdownlint-cli
	## Used by ReviewDog via .reviewdog.yml configs to lint markdown files in this project
	npm install markdownlint-cli

	## See https://golangci-lint.run/usage/install/#ci-installation
	## Used by ReviewDog
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh

	## See https://github.com/reviewdog/reviewdog#installation
	## Note: If yor're on Darwin with an ARM64 chip (Apple Silicon/Apple M1), know that
	##   at this point in time reviewdog doesn't publish Darmwin/arm64 binaries. You can
	##   get around this by either installing reviewdog manually or getting it via Homebrew
	##   and using Rosetta 2: https://support.apple.com/en-us/HT211861
ifeq (, $(shell command -v reviewdog --version))
	curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh
else
	@echo "Reviewdog installed at $$(which reviewdog)"
endif
.PHONY: linter-setup

lint-tf:
	terraform fmt -check -diff -recursive testdata

lint: lint-setup ## run lint andn print results
	./bin/reviewdog -conf .reviewdog.yml  -diff "git diff main"
.PHONY: lint

lint-ci: lint-setup ## run lint in CI, posting to PRs
	./bin/reviewdog -conf .reviewdog.yml  -reporter=github-pr-review -tee -level=info
.PHONY: lint-ci

lint-all: lint-setup ## run the fast go linters
	# doesn't seem to be a way to get reviewdog to not filter by diff
	./bin/reviewdog -conf .reviewdog.yml  -filter-mode nofilter
.PHONY: lint-all

TEMPLATES := $(shell find templates -not -name "*.go")

release: setup ## run a release
	./bin/bff bump
	git push
	goreleaser release
.PHONY: release

release-prerelease: setup ## release to github as a 'pre-release'
	go build ${LDFLAGS} .
	version=`./fogg version`; \
	git tag v"$$version"; \
	git push; \
	git push origin v"$$version";
	goreleaser release -f .goreleaser.prerelease.yml --debug
.PHONY: release-prerelease

release-snapshot: setup ## run a release
	goreleaser release --snapshot
.PHONY: release-snapshot

build: fmt ## build the binary
	go build ${LDFLAGS} .
.PHONY: build

coverage: ## run the go coverage tool, reading file coverage.out
	go tool cover -html=coverage.out
.PHONY: coverage

test: fmt ## run tests
 ifeq (, $(shell which gotest))
	go test -failfast -cover ./...
 else
	gotest -failfast -cover ./...
 endif
.PHONY: test

test-ci: ## run tests
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
.PHONY: test-ci

test-offline: ## run only tests that don't require internet
	go test -tags=offline ./...
.PHONY: test-offline

test-coverage: ## run the test with proper coverage reporting
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out
.PHONY: test-coverage

install: ## install the fogg binary in $GOPATH/bin
	go install ${LDFLAGS} .
.PHONY: install

help: ## display help for this makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

clean: ## clean the repo
	rm fogg 2>/dev/null || true
	go clean
	go clean -testcache
	rm -rf dist 2>/dev/null || true
	rm coverage.out 2>/dev/null || true

update-golden-files: clean ## update the golden files in testdata
	go test -v -run TestIntegration ./apply/ -update
.PHONY: update-golden-files
