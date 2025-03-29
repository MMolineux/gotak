.PHONY: all
all: buf lint test

.PHONY: buf
buf:
	buf generate

.PHONY: lint
lint:
	go vet ./...
	golangci-lint run

.PHONY: test
test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: clean
clean:
	rm -rf ./pkg/cotproto/*

.PHONY: changelog
changelog:
	changie batch auto
	changie merge

.PHONY: release
release:
	goreleaser release --clean

.PHONY: release-snapshot
release-snapshot:
	goreleaser release --snapshot --clean

.PHONY: setup
setup:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/miniscruff/changie@latest
	go install github.com/goreleaser/goreleaser@latest
