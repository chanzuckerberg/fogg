SHA=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)
DIRTY=$(shell if `git diff-index --quiet HEAD --`; then echo false; else echo true;  fi)
LDFLAGS=-ldflags "-w -s -X main.GitSha=${SHA} -X main.Version=${VERSION} -X main.Dirty=${DIRTY}"

build:
	@echo $(SHA)
	go build ${LDFLAGS} .

coverage:
	go tool cover -html=coverage.out

test:
	go test -cover ./...

install:
	go install .

.PHONY: build coverage test install