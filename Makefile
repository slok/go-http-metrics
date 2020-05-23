
UNIT_TEST_CMD := go test `go list ./... | grep -v vendor` -race
INTEGRATION_TEST_CMD := go test `go list ./... | grep -v vendor` -race -tags='integration'
BENCHMARK_CMD := go test `go list ./... | grep -v vendor` -benchmem -bench=.
CHECK_CMD = golangci-lint run -E goimports
DEPS_CMD := go mod tidy
MOCKS_CMD := go generate ./internal/mocks

.PHONY: default
default: test

.PHONY: unit-test
unit-test:
	$(UNIT_TEST_CMD)

.PHONY: integration-test
integration-test:
	$(INTEGRATION_TEST_CMD)

.PHONY: test
test: integration-test

.PHONY: benchmark
benchmark:
	$(BENCHMARK_CMD)

.PHONY: check
check: 
	$(CHECK_CMD)

.PHONY: deps
deps:
	$(DEPS_CMD)

.PHONY: mocks
mocks:
	$(MOCKS_CMD)

.PHONY: docs
docs: 
	godoc -http=":6060"
