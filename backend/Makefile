PROJECT_NAME=mmbrnr

MOCKS_DIR=./test/mocks
MOCKS_GEN_ARGS=--case underscore --keeptree --all --with-expecter --exported

ifeq ($(V),1)
	TEST_ARGS="-v"
endif

ifndef FILE_GLOB
	FILE_GLOB=./...
endif

.PHONY: build
build:
	go mod tidy
	go build -o /go/bin/$(PROJECT_NAME) main.go

.PHONY: serve
serve: build
	$(PROJECT_NAME) serve

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: gen-cmd-mocks
gen-cmd-mocks: 
	mockery $(MOCKS_GEN_ARGS) --output $(MOCKS_DIR)/cmd --outpkg cmdmocks --dir ./cmd

.PHONY: gen-internal-infra-mocks
gen-internal-infra-mocks:
	mockery --case underscore --with-expecter --exported --output $(MOCKS_DIR)/infra/rdb --outpkg rdbmocks --dir ./internal/infra/rdb --name connectionOpener
	mockery --case underscore --with-expecter --exported --output $(MOCKS_DIR)/infra/rdb/scm/postgres --outpkg scmpostgresmocks --dir ./internal/infra/rdb/scm/postgres --name connManager
	mockery --case underscore --with-expecter --exported --output $(MOCKS_DIR)/infra/rdb/schema --outpkg schema --dir ./internal/infra/rdb/schema --name connectionManager
	mockery --case underscore --with-expecter --exported --output $(MOCKS_DIR)/infra/rdb/schema --outpkg schema --dir ./internal/infra/rdb/schema --name hasSchema

.PHONY: gen-internal-pkg-mocks
gen-internal-pkg-mocks: 
	mockery --case underscore --with-expecter --exported --output $(MOCKS_DIR)/pkg/validation --outpkg validationmocks --dir ./internal/pkg/validation --name ValidationError

.PHONY: gen-mocks
gen-mocks: gen-cmd-mocks gen-internal-pkg-mocks gen-internal-infra-mocks

.PHONY: unit-test
unit-test: 
	go test -cover -coverprofile test/unit.out $(TEST_ARGS) $(FILE_GLOB)
	go tool cover -o=test/unit.html -html=test/unit.out

.PHONY: integration-test
integration-test:
	CI_INTEGRATION_TESTS_ENABLED=true go test -cover -coverprofile test/integration.out -p 1 -run "^TestIntegration" $(TEST_ARGS) $(FILE_GLOB)
	go tool cover -o=test/integration.html -html=test/integration.out

.PHONY: test
test: unit-test integration-test

.PHONY: ci
ci: fmt lint test