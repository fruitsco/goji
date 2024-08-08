.PHONY: all build clean test test-unit test-race test-msan staticcheck vet

all: build

# MARK: - Build

build:
	go build -o ./bin/roma $(GOFLAGS) .

# MARK: - Test

test: test-unit test-race

staticcheck:
	staticcheck ./...

vet:
	go vet ./...

test: test-unit

test-snap:
	UPDATE_SNAPS=true go test -covermode=count -coverprofile=coverage.out ./...

test-unit:
	go test -covermode=count -coverprofile=coverage.out ./...

test-race:
	go test -race ./...

test-msan:
	go test -msan ./...

# MARK: - Mockery

generate-mocks:
	mockery
