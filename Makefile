.PHONY: test
test:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: build
build:
	go build -o bin/pglint-cli cmd/pglint-cli/main.go

.PHONY:
docker:
	docker build -t pglint .