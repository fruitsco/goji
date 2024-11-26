.PHONY: all test test-unit test-race test-msan staticcheck vet

# MARK: - Test

test: test-unit test-race

staticcheck:
	staticcheck ./...

vet:
	go vet ./...

test: test-unit

test-unit:
	go test -covermode=count -coverprofile=coverage.out ./...

test-race:
	go test -race ./...

test-msan:
	go test -msan ./...
