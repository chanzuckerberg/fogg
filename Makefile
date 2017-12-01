
build:
	go build .

coverage:
	go tool cover -html=coverage.out

test:
	go test -coverprofile=coverage.out .

.PHONY: build coverage test