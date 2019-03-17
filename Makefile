
UNIT_TEST_CMD := go test `go list ./... | grep -v vendor` -race -v
INTEGRATION_TEST_CMD := go test `go list ./... | grep -v vendor`  -race -v -tags='integration'
BENCHMARK_CMD := go test -benchmem -bench=.
DEPS_CMD := GO111MODULE=on go mod tidy && GO111MODULE=on go mod vendor
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

.PHONY: ci
ci: test

.PHONY: deps
deps:
	$(DEPS_CMD)

.PHONY: mocks
mocks:
	$(MOCKS_CMD)

.PHONY: godoc
godoc: 
	godoc -http=":6060"
