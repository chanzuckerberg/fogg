
build:
	go build .

coverage:
	go tool cover -html=coverage.out

test:
	go test -cover ./...

.PHONY: build coverage test