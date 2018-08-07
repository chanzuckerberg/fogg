SHA=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)
DIRTY=$(shell if `git diff-index --quiet HEAD --`; then echo false; else echo true;  fi)
# TODO add release flag
LDFLAGS=-ldflags "-w -s -X github.com/chanzuckerberg/fogg/util.GitSha=${SHA} -X github.com/chanzuckerberg/fogg/util.Version=${VERSION} -X github.com/chanzuckerberg/fogg/util.Dirty=${DIRTY}"

all: test install

lint: ## run the fast go linters
	gometalinter --vendor --fast ./...

lint-slow: ## run all linters, even the slow ones
	gometalinter --vendor --deadline 120s ./...

packr: ## run the packr tool to generate our static files
	packr

release: packr
	goreleaser release --rm-dist

build: ## build the binary
	go build ${LDFLAGS} .

coverage: ## run the go coverage tool, reading file coverage.out
	go tool cover -html=coverage.out

test: ## run the tests
	go test -cover ./...

install: packr # install the fogg binary in $GOPATH/bin
	go install ${LDFLAGS} .

help: ## display help for this makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build coverage test install lint lint-slow packr release help
