
UNIT_TEST_CMD 			:= go test `go list ./... | grep -v test\/integration` -race -coverprofile=.test_coverage.txt && \
						   	go tool cover -func=.test_coverage.txt | tail -n1 | awk '{print "Total test coverage: " $$3}'
INTEGRATION_TEST_CMD 	:= go test ./test/integration -race
BENCHMARK_CMD 			:= go test `go list ./...` -benchmem -bench=.
CHECK_CMD 				:= golangci-lint run -E goimports
DEPS_CMD 				:= go mod tidy
MOCKS_CMD 				:= go generate ./internal/mocks

help: ## Show this help.
	@echo "Help"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[93m %s\n", $$1, $$2}'


.PHONY: default
default: test

.PHONY: unit-test
unit-test: ## Execute unit tests.
	$(UNIT_TEST_CMD)

.PHONY: integration-test
integration-test: ## Execute unit tests.
	$(INTEGRATION_TEST_CMD)

.PHONY: test ## Alias for unit tests.
test: unit-test

.PHONY: benchmark
benchmark: ## Execute benchmarks.
	$(BENCHMARK_CMD)

.PHONY: check
check: ## Execute check.
	$(CHECK_CMD)

.PHONY: deps
deps: ## Tidy dependencies.
	$(DEPS_CMD)

.PHONY: mocks
mocks: ## Generates mocks.
	$(MOCKS_CMD)

.PHONY: docs
docs: ## Runs docs example on :6060.
	godoc -http=":6060"
