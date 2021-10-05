ifeq ("$(origin V)", "command line")
  BUILD_VERBOSE = $(V)
endif
ifndef BUILD_VERBOSE
  BUILD_VERBOSE = 0
endif

ifeq ($(BUILD_VERBOSE),1)
  Q =
else
  Q = @
endif
ECHO := @echo

.PHONY: coverage
coverage: test
	$(Q)go tool cover -html=coverage.out -o coverage.html

.PHONY: test
test:
	$(Q)go test -v -covermode=count -coverprofile coverage.out ./...

.PHONY:
checkdeps:
	@echo "Checking dependencies"
ifeq ($(shell which golangci-lint),)
	$(ECHO) Installing golangci-lint
ifneq ($(shell which brew),)
	$(Q)brew install golangci-lint
else
	$(Q)mkdir -p ${GOPATH}/bin
	$(Q)curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.31.0
endif
endif

.PHONY: lint
lint: checkdeps
# see https://github.com/golangci/golangci-lint/issues/1040
	$(ECHO) "Running $@ check"
	@GO111MODULE=on $(shell which golangci-lint) cache clean
	@GO111MODULE=on $(shell which golangci-lint) run --timeout=5m --config ./.golangci.yml