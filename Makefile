SHA=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)
DIRTY=$(shell if `git diff-index --quiet HEAD --`; then echo false; else echo true;  fi)
# TODO add release flag
LDFLAGS=-ldflags "-w -s -X github.com/chanzuckerberg/fogg/util.GitSha=${SHA} -X github.com/chanzuckerberg/fogg/util.Version=${VERSION} -X github.com/chanzuckerberg/fogg/util.Dirty=${DIRTY}"

all: test install

lint:
	gometalinter --vendor --fast ./...

lint-slow:
	gometalinter --vendor --deadline 120s ./...

packr:
	packr

release: packr
	goreleaser release --rm-dist

build:
	go build ${LDFLAGS} .

coverage:
	go tool cover -html=coverage.out

test:
	go test -cover ./...

install:
	go install ${LDFLAGS} .

.PHONY: build coverage test install lint lint-slow packr release
