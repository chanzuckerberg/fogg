SHA=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)
DIRTY=false
# TODO add release flag
GO_PACKAGE=$(shell go list)
LDFLAGS=-ldflags "-w -s -X $(GO_PACKAGE)/util.GitSha=${SHA} -X $(GO_PACKAGE)/util.Version=${VERSION} -X $(GO_PACKAGE)/util.Dirty=${DIRTY}"
export GO111MODULE=on

all: test install

setup: ## setup development dependencies
	./.godownloader-packr.sh -d v1.24.1
	curl -sfL https://raw.githubusercontent.com/chanzuckerberg/bff/main/download.sh | sh
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh| sh
.PHONY: setup

fmt:
	goimports -w -d $$(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./dist/*")
.PHONY: fmt

lint: ## run lint andn print results
	./bin/reviewdog -conf .reviewdog.yml  -diff "git diff main" -tee
.PHONY: lint

lint-ci: ## run lint in CI, posting to PRs
	./bin/reviewdog -conf .reviewdog.yml  -reporter=github-pr-review -tee -level=info
.PHONY: lint-ci

lint-all: ## run the fast go linters
	# doesn't seem to be a way to get reviewdog to not filter by diff
	./bin/reviewdog -conf .reviewdog.yml  -filter-mode nofilter -tee
.PHONY: lint-all

TEMPLATES := $(shell find templates -not -name "*.go")

templates/a_templates-packr.go: $(TEMPLATES)
	./bin/packr clean -v
	./bin/packr -v

packr: templates/a_templates-packr.go ## run the packr tool to generate our static files
.PHONY: packr

docker: ## check to be sure docker is running
	@docker ps
.PHONY: docker

release: setup docker ## run a release
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

build: fmt packr ## build the binary
	go build ${LDFLAGS} .
.PHONY: build

deps:
	go mod tidy
.PHONY: deps

coverage: ## run the go coverage tool, reading file coverage.out
	go tool cover -html=coverage.out
.PHONY: coverage

test: fmt deps packr ## run tests
 ifeq (, $(shell which gotest))
	go test -failfast -cover ./...
 else
	gotest -failfast -cover ./...
 endif
.PHONY: test

test-ci: packr ## run tests
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
.PHONY: test-ci

test-offline: packr  ## run only tests that don't require internet
	go test -tags=offline ./...
.PHONY: test-offline

test-coverage: packr  ## run the test with proper coverage reporting
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out
.PHONY: test-coverage

install: packr ## install the fogg binary in $GOPATH/bin
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
	./bin/packr clean
	rm coverage.out 2>/dev/null || true

update-golden-files: clean deps ## update the golden files in testdata
	go test -v -run TestIntegration ./apply/ -update
.PHONY: update-golden-files
