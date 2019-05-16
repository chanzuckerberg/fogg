SHA=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)
DIRTY=$(shell if `git diff-index --quiet HEAD --`; then echo false; else echo true;  fi)
# TODO add release flag
LDFLAGS=-ldflags "-w -s -X github.com/chanzuckerberg/fogg/util.GitSha=${SHA} -X github.com/chanzuckerberg/fogg/util.Version=${VERSION} -X github.com/chanzuckerberg/fogg/util.Dirty=${DIRTY}"

all: test install

setup: ## setup development dependencies
	go get github.com/rakyll/gotest
	go install github.com/rakyll/gotest
	go get -u github.com/gobuffalo/packr/...
	go install github.com/gobuffalo/packr/packr
	curl -L https://raw.githubusercontent.com/chanzuckerberg/bff/master/download.sh | sh
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh

lint: ## run the fast go linters
	golangci-lint run
.PHONY: lint

packr: ## run the packr tool to generate our static files
	packr clean -v
	packr -v

release: ## run a release
	./bin/bff bump
	git push
	goreleaser release

release-prerelease: build ## release to github as a 'pre-release'
	version=`./fogg version`; \
	git tag v"$$version"; \
	git push
	git push --tags
	goreleaser release -f .goreleaser.prerelease.yml --debug

release-snapshot: ## run a release
	goreleaser release --snapshot

build: dep packr ## build the binary
	go build ${LDFLAGS} .

coverage: ## run the go coverage tool, reading file coverage.out
	go tool cover -html=coverage.out

test: dep ## run tests
	gotest -race -cover ./...

test-offline: dep ## run only tests that don't require internet
	gotest -tags=offline ./...
.PHONY: test-offline

test-coverage: ## run the test with proper coverage reporting
	goverage -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out
.PHONY: test-coverage

install: dep packr ## install the fogg binary in $GOPATH/bin
	go install ${LDFLAGS} .

help: ## display help for this makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## clean the repo
	rm fogg 2>/dev/null || true
	go clean 
	rm -rf dist 2>/dev/null || true
	packr clean
	rm coverage.out 2>/dev/null || true

dep: ## ensure dependencies are vendored
	dep ensure # this should be super-fast in the no-op case
.PHONY: dep

update-golden-files: clean dep ## update the golden files in testdata
	go test -v -run TestIntegration ./apply/ -update
.PHONY: update-golden-files

.PHONY: build clean coverage test install packr release help setup
