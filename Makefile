
build:
	go build .

coverage:
	go tool cover -html=coverage.out

test:
	go test -cover ./...

install:
	go install .

.PHONY: build coverage test